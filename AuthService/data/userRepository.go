package data

import (
	"github.com/DemoMicroservice/AuthService/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	C *mgo.Collection
}

func (r *UserRepository) GetAll(query *models.UserQuery) []model.Movie {
	string strQuery = "id"
	if query.Id != "" {
		r.C.Find(bson.M{"id":query.Id}).Iter
	}
}
func (r *UserRepository) Create(user *models.User) error {
	obj_id := bson.NewObjectId()
	user.Id = obj_id
	return r.C.Insert(&user)
}
