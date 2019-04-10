package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
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
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login LoginResource
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error in request body")
		return
	}
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("user")
	repo := &data.UserRepository{c}
	user := repo.Login(login.Username, login.Password)
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Wrong")
		return
	}
	common.GenerateToken(w, user)
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
		username, _ := claims["username"].(string)
		id, _ := claims["id"].(string)
		idEmployee, _ := claims["idEmployee"].(string)
		j, err := json.Marshal(AuthResource{
			Id:         id,
			IdEmployee: idEmployee,
			Username:   username,
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
