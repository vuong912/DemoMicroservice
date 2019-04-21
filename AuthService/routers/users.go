package routers

import (
	"net/http"

	"github.com/DemoMicroservice/AuthService/common"
	"github.com/DemoMicroservice/AuthService/controllers"
	"github.com/gorilla/mux"
)

func SetUsersRouters(router *mux.Router) *mux.Router {
	createUserHandler := http.HandlerFunc(controllers.CreateUserHandler)
	createUserRole := map[string]bool{
		common.AdminRole:   true,
		common.OwnerRole:   true,
		common.PlannerRole: true,
	}

	router.HandleFunc("/user/get", controllers.GetUsersHandler).Methods("GET")
	router.HandleFunc("/user/login", controllers.LoginHandler).Methods("POST")
	router.HandleFunc("/user/auth", controllers.AuthHandler).Methods("GET")
	router.Handle("/user/create", controllers.AuthMiddleware(createUserHandler, &createUserRole)).Methods("POST")

	return router
}
