package data

import (
	"github.com/DemoMicroservice/AuthService/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	C *mgo.Collection
}

func (r *UserRepository) GetAll(query interface{}, orderBy string, pageStep, pageSize int) (int, []models.User, error) {
	var users []models.User
	resQuery := r.C.Find(query)
	if orderBy != "" {
		resQuery = resQuery.Sort(orderBy)
	}
	sizeResult, err := resQuery.Count()
	if err != nil {
		return 0, nil, err
	}
	iter := resQuery.Skip((pageStep - 1) * pageSize).Limit(pageSize).Iter()
	err = iter.Err()
	if err != nil {
		return 0, nil, err
	}
	user := models.User{}
	for iter.Next(&user) {
		user.Password = ""
		users = append(users, user)
	}
	return sizeResult, users, nil
}
func (r *UserRepository) GetById(id string) (models.User, error) {
	var user models.User
	err := r.C.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)
	return user, err
}
func (r *UserRepository) Create(user *models.User) error {
	obj_id := bson.NewObjectId()
	user.Id = obj_id
	err := r.C.Insert(&user)
	user.Password = ""
	return err
}
func (r *UserRepository) Login(username, password string) (*models.User, error) {
	var user *models.User
	err := r.C.Find(bson.M{
		"username": username,
		"password": password,
	}).One(&user)
	user.Password = ""
	return user, err
}
