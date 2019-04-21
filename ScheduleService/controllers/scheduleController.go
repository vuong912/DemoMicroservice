package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/DemoMicroservice/ScheduleService/common"
	"github.com/DemoMicroservice/ScheduleService/data"
	"github.com/DemoMicroservice/ScheduleService/models"
	"gopkg.in/mgo.v2/bson"
)

func GetSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	var mapQuery = make(bson.M)

	vars := r.URL.Query()

	if val := vars.Get("id"); val != "" {
		mapQuery["_id"] = bson.ObjectIdHex(val)
	}

	dayMap := make(bson.M)
	layoutISO := "02-01-2006"
	if val := vars.Get("day"); val != "" {
		if value, err := time.Parse(layoutISO, val); err == nil {
			mapQuery["day"] = value
		}
	}
	if val := vars.Get("gteday"); val != "" {
		if gteValue, err := time.Parse(layoutISO, val); err == nil {
			dayMap["$gte"] = gteValue
			mapQuery["day"] = dayMap
		}

	}
	if val := vars.Get("lteday"); val != "" {
		if lteValue, err := time.Parse(layoutISO, val); err == nil {
			dayMap["$lte"] = lteValue
			mapQuery["day"] = dayMap
		}
	}

	if val := vars.Get("idshift"); val != "" {
		mapQuery["idShift"] = val
	}
	authInfo := tokenToInfo[r.Header.Get("Authorization")]

	if val := vars.Get("idbranch"); true {
		if authInfo.RoleName == common.PlannerRole && val != authInfo.IdBranch {
			common.DisplayAppError(w, errors.New("Not have permission."), "Not have permission.", http.StatusForbidden)
			return
		}
		if val != "" {
			mapQuery["idBranch"] = val
		}
	}
	if val := vars.Get("idemployee"); val != "" {
		mapQuery["detailSchedule.idEmployee"] = val
	}
	if val, err := strconv.ParseBool(vars.Get("check")); err == nil {
		mapQuery["detailSchedule.check"] = val
	}

	pageSize, err := strconv.Atoi(vars.Get("pagesize"))
	if err != nil {
		pageSize = 20
	}
	pageStep, err := strconv.Atoi(vars.Get("pagestep"))
	if err != nil {
		pageStep = 1
	}

	context := NewContext()
	defer context.Close()
	c := context.DbCollection("schedule")
	repo := &data.ScheduleRepository{c}
	size, schedules, err := repo.GetAll(mapQuery, vars.Get("orderby"), pageStep, pageSize)
	if err != nil {
		common.DisplayAppError(w, err, "Error query database", http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(SchedulesResource{
		Size: size,
		Data: schedules,
	})

	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
	//fmt.Printf("Result: %s\n", j)
}
func CreateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	var schedule models.Schedule
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		common.DisplayAppError(w, err, "Invalid employee data", http.StatusBadRequest)
		return
	}
	token := r.Header.Get("Authorization")
	authInfo := tokenToInfo[token]

	if authInfo.RoleName == common.PlannerRole && authInfo.IdBranch != schedule.IdBranch {
		common.DisplayAppError(w, err, "Not have permission in this branch", http.StatusForbidden)
	}

	schedule.CreatedBy = authInfo.IdEmployee
	schedule.ModifiedBy = authInfo.IdEmployee
	schedule.CreatedDay = time.Now()
	schedule.ModifiedDay = time.Now()

	context := NewContext()
	defer context.Close()
	c := context.DbCollection("employee")
	repo := &data.ScheduleRepository{c}
	err = repo.Create(&schedule)
	if err != nil {
		common.DisplayAppError(w, err, "Error create employee in database", http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(schedule)
	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
}
