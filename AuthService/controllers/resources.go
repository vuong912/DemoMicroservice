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
		IdUser     string `json:"idUser"`
		Username   string `json:"username"`
		IdEmployee string `json:"idEmployee"`
		Role       string `json:"role"`
	}
)
