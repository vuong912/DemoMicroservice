package data

import (
	"github.com/DemoMicroservice/ScheduleService/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ScheduleRepository struct {
	C *mgo.Collection
}

func (r *ScheduleRepository) GetAll(query interface{}, orderBy string, pageStep, pageSize int) (int, []models.Schedule, error) {
	var schedules []models.Schedule
	resQuery := r.C.Find(query)
	if orderBy != "" {
		resQuery = resQuery.Sort(orderBy)
	}
	sizeResult, err := resQuery.Count()
	if err != nil {
		return 0, nil, err
	}
	err = resQuery.Skip((pageStep - 1) * pageSize).Limit(pageSize).All(&schedules)
	return sizeResult, schedules, err
}

func (r *ScheduleRepository) Create(schedule *models.Schedule) error {
	scheduleExist := models.Schedule{}
	iter := r.C.Find(bson.M{
		"day":      schedule.Day,
		"idBranch": schedule.IdBranch,
		"idShift":  schedule.IdShift}).Iter()
	if err := iter.Err(); err != nil {
		return err
	}
	if iter.Next(&scheduleExist) {
		if schedule.DetailSchedule == nil || len(schedule.DetailSchedule) == 0 {
			err := r.C.Remove(bson.M{"_id": scheduleExist.Id})
			return err
		}
		filter := bson.M{"_id": scheduleExist.Id}
		update := bson.M{
			"$set": bson.M{
				"detailSchedule": schedule.DetailSchedule,
				"modifiedBy":     schedule.ModifiedBy,
				"modifiedDay":    schedule.ModifiedDay,
			},
		}
		return r.C.Update(filter, update)
	} else {
		if schedule.DetailSchedule == nil || len(schedule.DetailSchedule) == 0 {
			return nil
		}
		obj_id := bson.NewObjectId()
		schedule.Id = obj_id
		err := r.C.Insert(&schedule)
		return err
	}

}
