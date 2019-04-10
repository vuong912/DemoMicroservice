package controllers

import "github.com/DemoMicroservice/AuthService/models"

type (
	UsersResource struct {
		Size int           `json:"size"`
		Data []models.User `json:"data"`
	}
	LoginResource struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	AuthResource struct {
		Id         string `json:"id"`
		Username   string `json:"username"`
		IdEmployee string `json:"idEmployee"`
	}
)
