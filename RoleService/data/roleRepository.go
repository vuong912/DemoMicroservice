package data

import (
	"github.com/DemoMicroservice/RoleService/models"
	"gopkg.in/mgo.v2"
)

type RoleRepository struct {
	C *mgo.Collection
}

func (r *RoleRepository) GetAll(query interface{}, orderBy string, pageStep, pageSize int) (int, []models.Role) {
	var roles []models.Role
	resQuery := r.C.Find(query)
	if orderBy != "" {
		resQuery = resQuery.Sort(orderBy)
	}
	sizeResult, _ := resQuery.Count()
	resQuery.Skip((pageStep - 1) * pageSize).Limit(pageSize).All(&roles)
	return sizeResult, roles
}
