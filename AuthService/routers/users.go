package routers

import (
	"github.com/DemoMicroservice/AuthService/controllers"
	"github.com/gorilla/mux"
)

func SetUsersRouters(router *mux.Router) *mux.Router {
	router.HandleFunc("/user/get", controllers.GetUsersHandler).Methods("GET")
	router.HandleFunc("/user/login", controllers.LoginHandler).Methods("POST")
	router.HandleFunc("/user/auth", controllers.AuthHandler).Methods("GET")
	return router
}
