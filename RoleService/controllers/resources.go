package controllers

import "github.com/DemoMicroservice/RoleService/models"

type (
	RolesResource struct {
		Size int           `json:"size"`
		Data []models.Role `json:"data"`
	}
)
