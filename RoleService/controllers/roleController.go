package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DemoMicroservice/RoleService/common"
	"github.com/DemoMicroservice/RoleService/data"
	"gopkg.in/mgo.v2/bson"
)

func GetRolesHandler(w http.ResponseWriter, r *http.Request) {
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("role")
	repo := &data.RoleRepository{c}
	var mapQuery = make(bson.M)

	vars := r.URL.Query()

	if val := vars.Get("id"); val != "" {
		mapQuery["_id"] = bson.ObjectIdHex(val)
	}
	if val := vars.Get("rolename"); val != "" {
		mapQuery["roleName"] = bson.RegEx{val, ""}
	}
	if val := vars.Get("status"); val == "true" || val == "false" {
		mapQuery["status"] = (val == "true")
	}

	pageSize, err := strconv.Atoi(vars.Get("pagesize"))
	if err != nil {
		pageSize = 20
	}
	pageStep, err := strconv.Atoi(vars.Get("pagestep"))
	if err != nil {
		pageStep = 1
	}

	size, roles, err := repo.GetAll(mapQuery, vars.Get("orderby"), pageStep, pageSize)
	if err != nil {
		common.DisplayAppError(w, err, "Error query database", http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(RolesResource{
		Size: size,
		Data: roles,
	})

	if err != nil {
		common.DisplayAppError(w, err, "Error parse json", http.StatusInternalServerError)
		return
	}
	common.DisplayJsonResult(w, j)
	//fmt.Printf("Result: %s\n", j)
}
