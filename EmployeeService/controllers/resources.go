package controllers

import (
	"time"

	"github.com/DemoMicroservice/EmployeeService/models"
)

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
	Role struct {
		Id          string    `json:"id"`
		RoleName    string    `json:"roleName"`
		Status      bool      `json:"status"`
		CreatedDay  time.Time `json:"createdDay"`
		ModifiedDay time.Time `json:"modifiedDay"`
		CreatedBy   string    `json:"createdBy"`
		ModifiedBy  string    `json:"modifiedBy"`
	}
	RoleResource struct {
		Size int    `json"size"`
		Data []Role `json:"data"`
	}
	PermissionInfo struct {
		IdUser     string
		IdEmployee string
		RoleName   string
		IdBranch   string
	}
)
