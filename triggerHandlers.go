package main

import (
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/jakecoffman/gorunner/executor"
	"github.com/jakecoffman/gorunner/models"
	"github.com/martini-contrib/render"
)

func ListTriggers(r render.Render) {
	triggerList := models.GetTriggerList()

	r.JSON(200, triggerList.Json())
}

func AddTrigger(r render.Render, p martini.Params) {
	triggerList := models.GetTriggerList()

	trigger := models.Trigger{Name: p["name"]}
	triggerList.Append(trigger)
	r.JSON(200, trigger)
}

func GetTrigger(r render.Render, p martini.Params) {
	triggerList := models.GetTriggerList()

	trigger, err := triggerList.Get(p["trigger"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find trigger"})
		return
	}

	r.JSON(200, trigger)
}

func UpdateTrigger(r render.Render, p martini.Params, req *http.Request) {
	triggerList := models.GetTriggerList()

	trigger, err := triggerList.Get(p["trigger"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find trigger"})
		return
	}

	t := trigger.(models.Trigger)
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		r.JSON(400, Message{"message": "bad trigger payload"})
		return
	}
	t.Schedule = string(data)
	executor.AddTrigger(t)
	err = triggerList.Update(t)
	if err != nil {
		r.JSON(500, Message{"mesage": "unable to update trigger"})
	}
	r.JSON(200, t)
}

func DeleteTrigger(r render.Render, p martini.Params) {
	triggerList := models.GetTriggerList()

	triggerList.Delete(p["trigger"])

	r.JSON(200, Message{})
}

func ListJobsForTrigger(r render.Render, p martini.Params) {
	jobList := models.GetJobList()
	jobs := jobList.GetJobsWithTrigger(p["trigger"])
	r.JSON(200, jobs)
}
