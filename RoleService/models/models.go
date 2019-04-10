package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Role struct {
		Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
		RoleName    string        `bson:"roleName" json:"roleName"`
		Status      bool          `json:"status"`
		CreatedDay  time.Time     `bson:"createdDay" json:"createdDay"`
		ModifiedDay time.Time     `bson:"modifiedDay" json:"modifiedDay"`
		CreatedBy   string        `bson:"createdBy" json:"createdBy"`
		ModifiedBy  string        `bson:"modifiedBy" json:"modifiedBy"`
	}
)
