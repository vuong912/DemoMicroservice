package main

import (
	"log"
	"net/http"

	"github.com/DemoMicroservice/AuthService/common"
	"github.com/DemoMicroservice/AuthService/routers"
)

func main() {

	common.StartUp()
	router := routers.InitRoutes()

	server := &http.Server{
		Addr:    common.AppConfig.Server,
		Handler: router,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
