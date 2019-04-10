package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/DemoMicroservice/AuthService/models"
	"github.com/dgrijalva/jwt-go"
)

type (
	key struct {
		SecretKey []byte
	}
	Token struct {
		Token string `json:"token"`
	}
)

var Key key

func initKey() {
	data, err := ioutil.ReadFile("common/secret.key")
	if err != nil {
		fmt.Println(err)
	}
	Key = key{
		SecretKey: data,
	}
}

func GenerateToken(w http.ResponseWriter, user *models.User) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = user.Id
	claims["idEmployee"] = user.IdEmployee
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString(Key.SecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Printf("Error: %v\n", err)
		return
	}
	j, err := json.Marshal(Token{tokenString})
	if err != nil {
		log.Println("Error parse json key")
		return
	}
	DisplayJsonResult(w, j)
}
