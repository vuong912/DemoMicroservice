package controllers

import (
	"time"

	"github.com/DemoMicroservice/ScheduleService/models"
	"gopkg.in/mgo.v2/bson"
)

type (
	SchedulesResource struct {
		Size int               `json:"size"`
		Data []models.Schedule `json:"data"`
	}

	AuthResource struct {
		IdUser     string `json:"idUser"`
		Username   string `json:"username"`
		IdEmployee string `json:"idEmployee"`
		Role       string `json:"role"`
	}
	Role struct {
		Id          string    `json:"id"`
		RoleName    string    `json:"roleName"`
		Status      bool      `json:"status"`
		CreatedDay  time.Time `json:"createdDay"`
		ModifiedDay time.Time `json:"modifiedDay"`
		CreatedBy   string    `json:"createdBy"`
		ModifiedBy  string    `json:"modifiedBy"`
	}
	RoleResource struct {
		Size int    `json"size"`
		Data []Role `json:"data"`
	}
	PermissionInfo struct {
		IdUser     string
		IdEmployee string
		RoleName   string
		IdBranch   string
	}
	Employee struct {
		Id                 bson.ObjectId        `json:"id"`
		Name               string               `json:"name"`
		Birthday           time.Time            `json:"birthday"`
		PhoneNumber        string               `json:"phoneNumber"`
		Email              string               `json:"email"`
		Rangewage          float64              `json:"rangeWage"`
		IdBranch           string               `json:"idBranch"`
		Status             bool                 `json:"status"`
		DetailEmployeeWork []DetailEmployeeWork `json:"detailEmployeeWork"`
		CreatedDay         time.Time            `json:"createdDay"`
		ModifiedDay        time.Time            `json:"modifiedDay"`
		CreatedBy          string               `json:"createdBy"`
		ModifiedBy         string               `json:"modifiedBy"`
	}
	DetailEmployeeWork struct {
		IdShift string `json:"idShift"`
		Day     int    `json:"day"`
	}
)
