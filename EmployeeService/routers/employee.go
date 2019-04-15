package routers

import (
	"github.com/DemoMicroservice/EmployeeService/controllers"
	"github.com/gorilla/mux"
)

func SetUsersRouters(router *mux.Router) *mux.Router {
	router.HandleFunc("/employee/get", controllers.GetEmployeesHandler).Methods("GET")
	router.HandleFunc("/employee/getmyseft", controllers.GetMyseftHandler).Methods("GET")
	return router
}
