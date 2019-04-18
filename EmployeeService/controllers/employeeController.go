package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/DemoMicroservice/EmployeeService/common"
	"github.com/DemoMicroservice/EmployeeService/data"
	"github.com/DemoMicroservice/EmployeeService/models"
	"gopkg.in/mgo.v2/bson"
)

var tokenToInfo map[string]*PermissionInfo = make(map[string]*PermissionInfo)

func AuthMiddleware(next http.Handler, roleAccept *map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		authResource, err := GetAuthInfo(token)
		if err != nil {
			common.DisplayAppError(w, err, "Error authorization", http.StatusUnauthorized)
			return
		}
		role, err := GetRoleInfo(authResource.Role)
		if err != nil {
			common.DisplayAppError(w, err, "Error call api role", http.StatusInternalServerError)
			return
		}
		if val, ok := (*roleAccept)[role.RoleName]; !ok || val == false {
			common.DisplayAppError(w, errors.New("Unauthentication"), "Unauthentication", http.StatusForbidden)
			return
		}
		employee, err := GetEmployeeInfo(authResource.IdEmployee)
		if err != nil {
			common.DisplayAppError(w, err, "Error query database", http.StatusInternalServerError)
			return
		}

		permissionInfo := PermissionInfo{
			IdEmployee: authResource.IdEmployee,
			IdBranch:   employee.IdBranch,
			IdUser:     authResource.IdUser,
			RoleName:   role.RoleName,
		}
		tokenToInfo[token] = &permissionInfo
		next.ServeHTTP(w, r)
		delete(tokenToInfo, token)
	})
}
func GetRoleInfo(idRole string) (*Role, error) {
	bytes, err := common.RequestService(
		"GET",
		common.AppConfig.GetRoleAPIHost+"?id="+idRole,
		nil,
		"")
	if err != nil {
		return nil, err
	}
	roleresource := RoleResource{}
	err = json.Unmarshal(bytes, &roleresource)
	if err != nil {
		return nil, err
	}
	return &roleresource.Data[0], nil
}
func GetEmployeeInfo(idEmployee string) (*models.Employee, error) {
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")

	repo := &data.EmployeeRepository{c}
	employee, err := repo.GetById(idEmployee)
	return &employee, err
}
func GetAuthInfo(token string) (*AuthResource, error) {
	bytes, err := common.RequestService(
		"GET",
		common.AppConfig.AuthAPIHost,
		nil,
		token)
	if err != nil {
		return nil, err
	}
	auth := AuthResource{}
	err = json.Unmarshal(bytes, &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}
func GetMyseftHandler(w http.ResponseWriter, r *http.Request) {
	auth, err := GetAuthInfo(r.Header.Get("Authorization"))
	if err != nil {
		common.DisplayAppError(w, err, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//log.Println(*auth)
	//auth := AuthResource{IdEmployee: "5cb0372fe929393474fb7ff1"}
	employee, err := GetEmployeeInfo(auth.IdEmployee)
	if err != nil {
		common.DisplayAppError(w, err, "Error query database", http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(employee)

	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
}
func GetEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	var mapQuery = make(bson.M)

	vars := r.URL.Query()

	if val := vars.Get("id"); val != "" {
		mapQuery["_id"] = bson.ObjectIdHex(val)
	}
	if val := vars.Get("name"); val != "" {
		mapQuery["name"] = bson.RegEx{val, ""}
	}
	if val := vars.Get("phonenumber"); val != "" {
		mapQuery["phoneNumber"] = bson.RegEx{val, ""}
	}
	if val := vars.Get("email"); val != "" {
		mapQuery["email"] = bson.RegEx{val, ""}
	}

	rangeWageMap := make(bson.M)
	if val := vars.Get("gtrangewage"); val != "" {
		if gtValue, err := strconv.ParseFloat(val, 64); err == nil {
			rangeWageMap["$gt"] = gtValue
			mapQuery["rangeWage"] = rangeWageMap
		}

	}
	if val := vars.Get("ltrangewage"); val != "" {
		if ltValue, err := strconv.ParseFloat(val, 64); err == nil {
			rangeWageMap["$lt"] = ltValue
			mapQuery["rangeWage"] = rangeWageMap
		}
	}
	authInfo := tokenToInfo[r.Header.Get("Authorization")]
	if authInfo.RoleName == common.PlannerRole {
		mapQuery["idBranch"] = authInfo.IdBranch
	} else if val := vars.Get("idbranch"); val != "" {
		mapQuery["idBranch"] = val
	}
	if val := vars.Get("status"); val == "true" || val == "false" {
		mapQuery["status"] = (val == "true")
	}
	if val := vars.Get("idshift"); val != "" {
		mapQuery["detailEmployeeWork.idShift"] = val
	}
	if val, err := strconv.Atoi(vars.Get("day")); err == nil {
		mapQuery["detailEmployeeWork.day"] = val
	}

	pageSize, err := strconv.Atoi(vars.Get("pagesize"))
	if err != nil {
		pageSize = 20
	}
	pageStep, err := strconv.Atoi(vars.Get("pagestep"))
	if err != nil {
		pageStep = 1
	}

	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")
	repo := &data.EmployeeRepository{c}
	size, employees, err := repo.GetAll(mapQuery, vars.Get("orderby"), pageStep, pageSize)
	if err != nil {
		common.DisplayAppError(w, err, "Error query database", http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(EmployeesResource{
		Size: size,
		Data: employees,
	})

	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
	//fmt.Printf("Result: %s\n", j)
}
func CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid employee data", http.StatusBadRequest)
		return
	}
	employee.CreatedDay = time.Now()
	employee.ModifiedDay = time.Now()

	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")
	repo := &data.EmployeeRepository{c}
	err = repo.Create(&employee)
	if err != nil {
		common.DisplayAppError(w, err, "Error create database", http.StatusInternalServerError)
	}
	j, err := json.Marshal(employee)
	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
}
