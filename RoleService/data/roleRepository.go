package data

import (
	"github.com/DemoMicroservice/RoleService/models"
	"gopkg.in/mgo.v2"
)

type RoleRepository struct {
	C *mgo.Collection
}

func (r *RoleRepository) GetAll(query interface{}, orderBy string, pageStep, pageSize int) (int, []models.Role, error) {
	var roles []models.Role
	resQuery := r.C.Find(query)
	if orderBy != "" {
		resQuery = resQuery.Sort(orderBy)
	}
	sizeResult, err := resQuery.Count()
	if err != nil {
		return 0, nil, err
	}
	err = resQuery.Skip((pageStep - 1) * pageSize).Limit(pageSize).All(&roles)
	return sizeResult, roles, err
}
