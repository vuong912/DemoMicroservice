package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DemoMicroservice/RoleService/common"
	"github.com/DemoMicroservice/RoleService/routers"
	"github.com/gorilla/handlers"
)

func main() {

	common.StartUp()
	router := routers.InitRoutes()
	logFile, err := os.OpenFile("role_server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(
		common.AppConfig.Server,
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(handlers.LoggingHandler(logFile, router))))

	/*server := &http.Server{
		Addr:    common.AppConfig.Server,
		Handler: router,
	}
	server.ListenAndServe()*/
}
