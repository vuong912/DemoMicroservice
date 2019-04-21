package routers

import (
	"net/http"

	"github.com/DemoMicroservice/ScheduleService/common"
	"github.com/DemoMicroservice/ScheduleService/controllers"
	"github.com/gorilla/mux"
)

func SetUsersRouters(router *mux.Router) *mux.Router {
	//router.HandleFunc("/employee/get", controllers.GetEmployeesHandler).Methods("GET")
	getSchedulesHandler := http.HandlerFunc(controllers.GetSchedulesHandler)
	getSchedulesRole := map[string]bool{
		common.AdminRole:   true,
		common.OwnerRole:   true,
		common.PlannerRole: true,
	}

	createScheduleHandler := http.HandlerFunc(controllers.CreateScheduleHandler)
	createScheduleRole := map[string]bool{
		common.AdminRole:   true,
		common.OwnerRole:   true,
		common.PlannerRole: true,
	}

	router.Handle("/schedule/get", controllers.AuthMiddleware(getSchedulesHandler, &getSchedulesRole)).Methods("GET")
	router.Handle("/schedule/create", controllers.AuthMiddleware(createScheduleHandler, &createScheduleRole)).Methods("POST")
	return router
}
