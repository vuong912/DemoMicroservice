package data

import (
	"github.com/DemoMicroservice/EmployeeService/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type EmployeeRepository struct {
	C *mgo.Collection
}

func (r *EmployeeRepository) GetAll(query interface{}, orderBy string, pageStep, pageSize int) (int, []models.Employee) {
	var employees []models.Employee
	resQuery := r.C.Find(query)
	if orderBy != "" {
		resQuery = resQuery.Sort(orderBy)
	}
	sizeResult, _ := resQuery.Count()
	resQuery.Skip((pageStep - 1) * pageSize).Limit(pageSize).All(&employees)
	return sizeResult, employees
}
func (r *EmployeeRepository) GetById(id string) models.Employee {
	var employee models.Employee
	r.C.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&employee)
	return employee
}
