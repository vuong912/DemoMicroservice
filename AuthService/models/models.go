package models

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	User struct {
		Id         bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Username   string        `json:"username"`
		Password   string        `json:"password"`
		IdEmployee string        `bson:"idEmployee" json:"idEmployee"`
		Role       string        `json:"role"`
	}
)
