package data

import (
	"time"

	"github.com/DemoMicroservice/EmployeeService/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type EmployeeRepository struct {
	C *mgo.Collection
}

func (r *EmployeeRepository) GetAll(query interface{}, orderBy string, pageStep, pageSize int) (int, []models.Employee, error) {
	var employees []models.Employee
	resQuery := r.C.Find(query)
	if orderBy != "" {
		resQuery = resQuery.Sort(orderBy)
	}
	sizeResult, err := resQuery.Count()
	if err != nil {
		return 0, nil, err
	}
	err = resQuery.Skip((pageStep - 1) * pageSize).Limit(pageSize).All(&employees)
	return sizeResult, employees, err
}
func (r *EmployeeRepository) GetById(id string) (models.Employee, error) {
	var employee models.Employee
	err := r.C.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&employee)
	return employee, err
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
	obj_id := bson.NewObjectId()
	employee.Id = obj_id
	err := r.C.Insert(&employee)
	return err
}
func (r *EmployeeRepository) Update(id, column, modifiedBy string, data interface{}) error {
	filter := bson.M{"_id": bson.ObjectIdHex(id)}
	update := bson.M{
		"$set": bson.M{
			column:        data,
			"modifiedBy":  modifiedBy,
			"modifiedDay": time.Now(),
		},
	}
	return r.C.Update(filter, update)
}
