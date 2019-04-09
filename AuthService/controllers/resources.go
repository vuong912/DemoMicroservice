package controllers

import "github.com/DemoMicroservice/AuthService/models"

type (
	UsersResource struct {
		Size int           `json:"size"`
		Data []models.User `json:"data"`
	}
)
