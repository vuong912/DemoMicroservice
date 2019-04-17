package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DemoMicroservice/EmployeeService/common"
	"github.com/DemoMicroservice/EmployeeService/data"
	"github.com/DemoMicroservice/EmployeeService/models"
	"gopkg.in/mgo.v2/bson"
)

func GetAuthInfo(r *http.Request) (*AuthResource, error) {
	bytes, err := common.RequestService(
		"GET",
		common.AppConfig.AuthAPIHost,
		nil,
		"",
		r.Header.Get("Authorization"))
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
	auth, err := GetAuthInfo(r)
	if err != nil {
		common.DisplayAppError(w, err, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(*auth)
	//auth := AuthResource{IdEmployee: "5cb0372fe929393474fb7ff1"}
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")

	repo := &data.EmployeeRepository{c}
	employee, err := repo.GetById((*auth).IdEmployee)
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
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")
	repo := &data.EmployeeRepository{c}
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

	if val := vars.Get("idbranch"); val != "" {
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
	// Decode the incoming Movie json
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid employee data", http.StatusInternalServerError)
		return
	}

	// create new context
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")
	// Insert a movie document
	repo := &data.EmployeeRepository{c}
	employee.CreatedDay = time.Now()
	employee.ModifiedDay = time.Now()
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
