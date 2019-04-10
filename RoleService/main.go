package main

import (
	"log"
	"net/http"

	"github.com/DemoMicroservice/RoleService/common"
	"github.com/DemoMicroservice/RoleService/routers"
	"github.com/gorilla/handlers"
)

func main() {

	common.StartUp()
	router := routers.InitRoutes()
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(
		common.AppConfig.Server,
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router)))
}
