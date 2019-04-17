package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func GenerateToken(user *models.User) ([]byte, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = user.Id
	claims["idEmployee"] = user.IdEmployee
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString(Key.SecretKey)
	if err != nil {
		return nil, err
	}
	j, err := json.Marshal(Token{tokenString})
	if err != nil {
		return nil, err
	}
	return j, nil
}
