package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/DemoMicroservice/AuthService/common"
	"github.com/DemoMicroservice/AuthService/data"
	"github.com/dgrijalva/jwt-go"
)

var tokenToInfo map[string]*PermissionInfo = make(map[string]*PermissionInfo)

func AuthMiddleware(next http.Handler, roleAccept *map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		authResource, err := GetAuthInfo(token)
		if err != nil {
			common.DisplayAppError(w, err, "Error authorization", http.StatusUnauthorized)
			return
		}
		role, err := GetRoleInfo(authResource.Role)
		if err != nil {
			common.DisplayAppError(w, err, "Error call api role", http.StatusInternalServerError)
			return
		}
		if val, ok := (*roleAccept)[role.RoleName]; !ok || val == false {
			common.DisplayAppError(w, errors.New("Unauthentication"), "Unauthentication", http.StatusForbidden)
			return
		}
		employee, err := GetEmployeeInfo(authResource.IdEmployee, token)
		if err != nil {
			common.DisplayAppError(w, err, "Error query database", http.StatusInternalServerError)
			return
		}
		if employee.Status == false {
			common.DisplayAppError(w, errors.New("Not permission"), "Your account have been blocked.", http.StatusForbidden)
			return
		}
		permissionInfo := PermissionInfo{
			IdEmployee: authResource.IdEmployee,
			IdBranch:   employee.IdBranch,
			IdUser:     authResource.IdUser,
			RoleName:   role.RoleName,
		}
		tokenToInfo[token] = &permissionInfo
		next.ServeHTTP(w, r)
		delete(tokenToInfo, token)
	})
}
func GetRoleInfo(idRole string) (*Role, error) {
	bytes, err := common.RequestService(
		"GET",
		common.AppConfig.GetRoleAPIHost+"?id="+idRole,
		nil,
		"")
	if err != nil {
		return nil, err
	}
	roleresource := RoleResource{}
	err = json.Unmarshal(bytes, &roleresource)
	if err != nil {
		return nil, err
	}
	return &roleresource.Data[0], nil
}
func GetEmployeeInfo(idEmployee, token string) (*Employee, error) {
	bytes, err := common.RequestService(
		"GET",
		common.AppConfig.GetMySelfEmployeeAPIHost,
		nil,
		token)
	if err != nil {
		return nil, err
	}
	employee := Employee{}
	err = json.Unmarshal(bytes, &employee)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func GetAuthInfo(token string) (*AuthResource, error) {
	if token == "" {
		return nil, errors.New("Unauthorized")
	}
	token = strings.TrimPrefix(token, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, (errors.New("Unexpected signing method"))
		}
		return common.Key.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, errors.New("Unauthorized")
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	id, _ := claims["id"].(string)
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("user")
	repo := &data.UserRepository{c}
	user, err := repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return &AuthResource{
		IdUser:     id,
		IdEmployee: user.IdEmployee,
		Username:   user.Username,
		Role:       user.Role,
	}, nil
}
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	authResource, err := GetAuthInfo(token)
	if err != nil {
		common.DisplayAppError(w, err, "Error when authorization", http.StatusUnauthorized)
		return
	}
	j, err := json.Marshal(authResource)
	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
}
