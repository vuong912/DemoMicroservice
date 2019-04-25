package routers

import (
	"net/http"

	"github.com/DemoMicroservice/EmployeeService/common"

	"github.com/DemoMicroservice/EmployeeService/controllers"
	"github.com/gorilla/mux"
)

func SetUsersRouters(router *mux.Router) *mux.Router {
	//router.HandleFunc("/employee/get", controllers.GetEmployeesHandler).Methods("GET")
	getEmployeesHandler := http.HandlerFunc(controllers.GetEmployeesHandler)
	getEmployeesRole := map[string]bool{
		common.AdminRole:   true,
		common.OwnerRole:   true,
		common.PlannerRole: true,
	}

	createEmployeeHandler := http.HandlerFunc(controllers.CreateEmployeeHandler)
	createEmployeeRole := map[string]bool{
		common.AdminRole:   true,
		common.OwnerRole:   true,
		common.PlannerRole: true,
	}

	updateEmployeeWorkHandler := http.HandlerFunc(controllers.UpdateEmployeeWorkHandler)
	updateEmployeeWorkRole := map[string]bool{
		common.AdminRole:   true,
		common.OwnerRole:   true,
		common.PlannerRole: true,
	}

	router.Handle("/get", controllers.AuthMiddleware(getEmployeesHandler, &getEmployeesRole)).Methods("GET")
	router.HandleFunc("/getmyself", controllers.GetMyseftHandler).Methods("GET")
	router.Handle("/create", controllers.AuthMiddleware(createEmployeeHandler, &createEmployeeRole)).Methods("POST")
	router.Handle("/update/work", controllers.AuthMiddleware(updateEmployeeWorkHandler, &updateEmployeeWorkRole)).Methods("POST")

	return router
}
