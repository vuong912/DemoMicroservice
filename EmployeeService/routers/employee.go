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
	router.Handle("/employee/get", controllers.AuthMiddleware(getEmployeesHandler, &getEmployeesRole))
	router.HandleFunc("/employee/getmyseft", controllers.GetMyseftHandler).Methods("GET")
	return router
}
