package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Employee struct {
		Id                 bson.ObjectId        `bson:"_id,omitempty" json:"id"`
		Name               string               `bson:"name" json:"name"`
		Birthday           time.Time            `bson:"birthday" json:"birthday"`
		PhoneNumber        string               `bson:"phoneNumber" json:"phoneNumber"`
		Email              string               `bson:"email" json:"email"`
		Rangewage          float64              `bson:"rangeWage" json:"rangeWage"`
		IdBranch           string               `bson:"idBranch" json:"idBranch"`
		Status             bool                 `bson:"status" json:"status"`
		DetailEmployeeWork []DetailEmployeeWork `bson:"detailEmployeeWork" json:"detailEmployeeWork"`
		CreatedDay         time.Time            `bson:"createdDay" json:"createdDay"`
		ModifiedDay        time.Time            `bson:"modifiedDay" json:"modifiedDay"`
		CreatedBy          string               `bson:"createdBy" json:"createdBy"`
		ModifiedBy         string               `bson:"modifiedBy" json:"modifiedBy"`
	}
	DetailEmployeeWork struct {
		IdShift string `bson:"idShift" json:"idShift"`
		Day     int    `bson:"day" json:"day"`
	}
)
