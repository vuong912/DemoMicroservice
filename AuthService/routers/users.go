package routers

import (
	"github.com/DemoMicroservice/AuthService/controllers"
	"github.com/gorilla/mux"
)

func SetUsersRouters(router *mux.Router) *mux.Router {
	router.HandleFunc("/users", controllers.GetUsersHandler).Methods("GET")
	return router
}
