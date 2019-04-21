package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DemoMicroservice/AuthService/common"
	"github.com/DemoMicroservice/AuthService/data"
	"github.com/DemoMicroservice/AuthService/models"
	"gopkg.in/mgo.v2/bson"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("user")
	repo := &data.UserRepository{c}
	var mapQuery = make(bson.M)
	vars := r.URL.Query()

	if val := vars.Get("id"); val != "" {
		mapQuery["_id"] = bson.ObjectIdHex(val)
	}
	if val := vars.Get("username"); val != "" {
		mapQuery["username"] = bson.RegEx{val, ""}
	}
	if val := vars.Get("idemployee"); val != "" {
		mapQuery["idEmployee"] = val
	}
	if val := vars.Get("role"); val != "" {
		mapQuery["role"] = val
	}

	pageSize, err := strconv.Atoi(vars.Get("pagesize"))
	if err != nil {
		pageSize = 20
	}
	pageStep, err := strconv.Atoi(vars.Get("pagestep"))
	if err != nil {
		pageStep = 1
	}

	size, users, err := repo.GetAll(mapQuery, vars.Get("orderby"), pageStep, pageSize)
	if err != nil {
		common.DisplayAppError(w, err, "Error query database", http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(UsersResource{
		Size: size,
		Data: users,
	})
	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
	//fmt.Printf("Result: %s\n", j)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserResource CreateUserResource
	err := json.NewDecoder(r.Body).Decode(&createUserResource)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid user data", http.StatusBadRequest)
		return
	}
	user := models.User{
		Username:   createUserResource.Username,
		Password:   common.GetMD5Hash(common.AppConfig.DefaultPassword),
		Role:       createUserResource.Role,
		IdEmployee: createUserResource.IdEmployee,
	}
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("user")
	repo := &data.UserRepository{c}
	err = repo.Create(&user)
	if err != nil {
		common.DisplayAppError(w, err, "Error when create user", http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(user)
	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login LoginResource
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		common.DisplayAppError(w, err, "Failed parse json", http.StatusBadRequest)
		return
	}
	login.Password = common.GetMD5Hash(login.Password)
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("user")
	repo := &data.UserRepository{c}
	user, err := repo.Login(login.Username, login.Password)
	if err != nil || user == nil {
		common.DisplayAppError(w, err, "Failed login", http.StatusUnauthorized)
		return
	}
	bytes, err := common.GenerateToken(user)
	if err != nil || bytes == nil {
		common.DisplayAppError(w, err, "Failed generate token", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, bytes)
}
