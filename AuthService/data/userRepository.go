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
	iter := resQuery.Skip((pageStep - 1) * pageSize).Limit(pageSize).Iter()
	user := models.User{}
	for iter.Next(&user) {
		user.Password = ""
		users = append(users, user)
	}
	return sizeResult, users
}
func (r *UserRepository) GetById(id string) models.User {
	var user models.User
	r.C.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)
	return user
}
func (r *UserRepository) Create(user *models.User) error {
	obj_id := bson.NewObjectId()
	user.Id = obj_id
	return r.C.Insert(&user)
}
func (r *UserRepository) Login(username, password string) *models.User {
	var user *models.User
	r.C.Find(bson.M{
		"username": username,
		"password": password,
	}).One(&user)
	return user
}
