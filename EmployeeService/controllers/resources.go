package controllers

import "github.com/DemoMicroservice/EmployeeService/models"

type (
	EmployeesResource struct {
		Size int               `json:"size"`
		Data []models.Employee `json:"data"`
	}
	AuthResource struct {
		IdUser     string `json:"idUser"`
		Username   string `json:"username"`
		IdEmployee string `json:"idEmployee"`
		Role       string `json:"role"`
	}
)
