package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/DemoMicroservice/AuthService/common"
	"github.com/DemoMicroservice/AuthService/data"
	"gopkg.in/mgo.v2/bson"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("user")
	repo := &data.UserRepository{c}
	var mapQuery = make(bson.M)

	vars := r.URL.Query()

	pageSize, err := strconv.Atoi(vars.Get("pagesize"))
	if err != nil {
		pageSize = 20
	}
	pageStep, err := strconv.Atoi(vars.Get("pagestep"))
	if err != nil {
		pageStep = 1
	}
	log.Printf("%d %d\n", pageSize, pageStep)
	size, users := repo.GetAll(mapQuery, vars.Get("orderby"), pageStep, pageSize)
	j, err := json.Marshal(UsersResource{
		Size: size,
		Data: users,
	})

	if err != nil {
		log.Println("Error parse json")
		return
	}
	common.DisplayJsonResult(w, j)
	//fmt.Printf("Result: %s\n", j)
}
