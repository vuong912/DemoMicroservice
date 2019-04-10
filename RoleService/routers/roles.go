package routers

import (
	"github.com/DemoMicroservice/RoleService/controllers"
	"github.com/gorilla/mux"
)

func SetUsersRouters(router *mux.Router) *mux.Router {
	router.HandleFunc("/roles", controllers.GetRolesHandler).Methods("GET")
	return router
}
