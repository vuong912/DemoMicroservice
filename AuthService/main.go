package main

import (
	"log"
	"net/http"

	"github.com/DemoMicroservice/demo_microservice/Template/common"
	"github.com/DemoMicroservice/demo_microservice/Template/routers"
)

func main() {
	ConfigTest()
	common.StartUp()
	router := routers.InitRoutes()
	server := &http.Server{
		Addr:    common.AppConfig.Server,
		Handler: router,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
