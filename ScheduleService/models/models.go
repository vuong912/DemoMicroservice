package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Schedule struct {
		Id             bson.ObjectId    `bson:"_id,omitempty" json:"id"`
		Day            time.Time        `bson:"day" json:"day"`
		IdShift        string           `bson:"idShift" json:"idShift"`
		IdBranch       string           `bson:"idBranch" json:"idBranch"`
		DetailSchedule []DetailSchedule `bson:"detailSchedule" json:"detailSchedule"`
		CreatedDay     time.Time        `bson:"createdDay" json:"createdDay"`
		ModifiedDay    time.Time        `bson:"modifiedDay" json:"modifiedDay"`
		CreatedBy      string           `bson:"createdBy" json:"createdBy"`
		ModifiedBy     string           `bson:"modifiedBy" json:"modifiedBy"`
	}
	DetailSchedule struct {
		IdEmployee string `bson:"idEmployee" json:"idEmployee"`
		Check      bool   `bson:"check" json:"check"`
	}
)
