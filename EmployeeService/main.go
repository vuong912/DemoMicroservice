package main

import (
	"log"
	"net/http"

	"github.com/DemoMicroservice/EmployeeService/common"
	"github.com/DemoMicroservice/EmployeeService/routers"
	"github.com/gorilla/handlers"
)

func main() {

	common.StartUp()
	router := routers.InitRoutes()
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(
		common.AppConfig.Server,
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router)))
	/*
		server := &http.Server{
			Addr:    common.AppConfig.Server,
			Handler: router,
		}
		server.ListenAndServe()*/
}