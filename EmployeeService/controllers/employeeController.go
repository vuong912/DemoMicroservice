package controllers

import (
	"bytes"
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
	if val := vars.Get("gterangewage"); val != "" {
		if gtValue, err := strconv.ParseFloat(val, 64); err == nil {
			rangeWageMap["$gte"] = gtValue
			mapQuery["rangeWage"] = rangeWageMap
		}

	}
	if val := vars.Get("lterangewage"); val != "" {
		if ltValue, err := strconv.ParseFloat(val, 64); err == nil {
			rangeWageMap["$lte"] = ltValue
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
	err = ValidateEmployee(&employee)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid employee data", http.StatusBadRequest)
		return
	}
	token := r.Header.Get("Authorization")
	authInfo := tokenToInfo[token]

	if authInfo.RoleName == common.PlannerRole && authInfo.IdBranch != employee.IdBranch {
		common.DisplayAppError(w, err, "Not have permission in this branch", http.StatusForbidden)
	}

	employee.CreatedBy = authInfo.IdEmployee
	employee.ModifiedBy = authInfo.IdEmployee
	employee.CreatedDay = time.Now()
	employee.ModifiedDay = time.Now()

	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")
	repo := &data.EmployeeRepository{c}
	err = repo.Create(&employee)
	if err != nil {
		common.DisplayAppError(w, err, "Error create employee in database", http.StatusInternalServerError)
		return
	}
	role, err := GetRoleByName(common.EmployeeRole)
	if err != nil {
		common.DisplayAppError(w, err, "Error get role", http.StatusInternalServerError)
		return
	}
	createUserResource := CreateUserResource{
		Username:   employee.Email,
		IdEmployee: employee.Id.Hex(),
		Role:       role.Id,
	}
	dataCreateUser, err := json.Marshal(createUserResource)
	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	err = CreateUser(dataCreateUser, token)
	if err != nil {
		common.DisplayAppError(w, err, "Error create user", http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(employee)
	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
}
func GetRoleByName(name string) (*Role, error) {
	bytes, err := common.RequestService(
		"GET",
		common.AppConfig.GetRoleAPIHost+"?rolename="+name,
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
func CreateUser(data []byte, token string) error {
	_, err := common.RequestService("POST",
		common.AppConfig.CreateUserAPIHost,
		bytes.NewBuffer(data),
		token)
	return err
}

func UpdateEmployeeWorkHandler(w http.ResponseWriter, r *http.Request) {
	var employeeUpdateWorkResource EmployeeUpdateWorkResource
	err := json.NewDecoder(r.Body).Decode(&employeeUpdateWorkResource)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid employee data", http.StatusBadRequest)
		return
	}
	authInfo := tokenToInfo[r.Header.Get("Authorization")]

	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")
	repo := &data.EmployeeRepository{c}
	if authInfo.RoleName == common.PlannerRole {
		size, employees, err := repo.GetAll(bson.M{"_id": employeeUpdateWorkResource.IdEmployee}, "", 1, 1)
		if err != nil {
			common.DisplayAppError(w, err, "Error query database", http.StatusInternalServerError)
			return
		}
		if size == 0 {
			common.DisplayAppError(w, errors.New("Wrong id"), "Not find this employee", http.StatusBadRequest)
			return
		}
		if authInfo.IdBranch != employees[0].IdBranch {
			common.DisplayAppError(w, errors.New("Not have permission"), "Not have permission", http.StatusForbidden)
			return
		}
	}
	err = repo.Update(employeeUpdateWorkResource.IdEmployee,
		"detailEmployeeWork",
		authInfo.IdEmployee,
		employeeUpdateWorkResource.DetailEmployeeWork)
	if err != nil {
		common.DisplayAppError(w, err, "Update fail", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, nil)

}
