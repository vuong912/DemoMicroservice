package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DemoMicroservice/demo_microservice/Template/common"
	"github.com/DemoMicroservice/demo_microservice/Template/routers"
)

func ConfigTest() {
	fmt.Println("Test...")
	fmt.Println(common.AppConfig.Server)
	fmt.Println(common.AppConfig.ServerIP)
	fmt.Println(common.AppConfig.ServiceHost)
	fmt.Println(common.AppConfig.DBUser)
	fmt.Println(common.AppConfig.DBPwd)
	fmt.Println(common.AppConfig.Database)
	fmt.Println(common.AppConfig.MongoDBHost)
}

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
