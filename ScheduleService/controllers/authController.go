package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DemoMicroservice/ScheduleService/common"
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
	bytes, err := common.RequestService(
		"GET",
		common.AppConfig.AuthAPIHost,
		nil,
		token)
	if err != nil {
		return nil, err
	}
	auth := AuthResource{}
	err = json.Unmarshal(bytes, &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}
