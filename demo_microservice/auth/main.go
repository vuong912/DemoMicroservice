package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

type (
	User struct {
		Username string `json:"username"`
		Password string `json:"password`
	}
	Token struct {
		Token string `json:"token"`
	}
	Account struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}
)

var secretKey = "secret123aAabAB"

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error in request body")
		return
	}
	if user.Username != "admin" || user.Password != "123456" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Wrong")
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["iss"] = "admin"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Printf("Error: %v\n", err)
		return
	}
	jsonResponse(Token{tokenString}, w)

}
func jsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func authHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var token string
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, (errors.New("Unexpected signing method"))
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, err
	}
	if parsedToken.Valid {
		w.WriteHeader(http.StatusOK)
		claims := parsedToken.Claims.(jwt.MapClaims)
		username, _ := claims["iss"].(string)
		return Account{
			Name: username + " AA",
			Role: "AD",
		}, nil
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, errors.New("expire")
	}
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	account, err := authHandler(w, r)
	if err != nil {
		fmt.Fprintln(w, "Unauthorized")
		log.Println(err)
		return
	}
	jsonResponse(account, w)
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/", indexHandler).Methods("GET")
	http.ListenAndServe(":8080", router)
}
