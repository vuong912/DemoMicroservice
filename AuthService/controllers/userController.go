package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DemoMicroservice/AuthService/common"
	"github.com/DemoMicroservice/AuthService/data"
	"github.com/dgrijalva/jwt-go"
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
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login LoginResource
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		common.DisplayAppError(w, err, "Failed parse json", http.StatusBadRequest)
		return
	}
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
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var token string
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, (errors.New("Unexpected signing method"))
		}
		return common.Key.SecretKey, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if parsedToken.Valid {
		claims := parsedToken.Claims.(jwt.MapClaims)
		id, _ := claims["id"].(string)
		//username, _ := claims["username"].(string)
		//idEmployee, _ := claims["idEmployee"].(string)

		context := NewContext()
		defer context.Close()
		c := context.DbCollection("user")
		repo := &data.UserRepository{c}
		user, err := repo.GetById(id)
		if err != nil {
			common.DisplayAppError(w, err, "Error when get user by id", http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(AuthResource{
			IdUser:     id,
			IdEmployee: user.IdEmployee,
			Username:   user.Username,
			Role:       user.Role,
		})

		if err != nil {
			log.Println("Error parse json")
			return
		}
		common.DisplayJsonResult(w, j)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
