package data

import (
	"github.com/DemoMicroservice/AuthService/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	C *mgo.Collection
}

func (r *UserRepository) GetAll(query interface{}, orderBy string, pageStep, pageSize int) (int, []models.User) {
	var users []models.User
	resQuery := r.C.Find(query)
	if orderBy != "" {
		resQuery = resQuery.Sort(orderBy)
	}
	sizeResult, _ := resQuery.Count()
	resQuery.Skip((pageStep - 1) * pageSize).Limit(pageSize).All(&users)
	return sizeResult, users
}
func (r *UserRepository) Create(user *models.User) error {
	obj_id := bson.NewObjectId()
	user.Id = obj_id
	return r.C.Insert(&user)
}
