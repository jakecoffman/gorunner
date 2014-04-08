package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
	"github.com/jakecoffman/gorunner/models"
	"github.com/martini-contrib/render"
	"github.com/nu7hatch/gouuid"
)

func ListRuns(r render.Render, req *http.Request) {
	runsList := models.GetRunListSorted()

	offset := req.FormValue("offset")
	length := req.FormValue("length")

	if offset == "" {
		offset = "-1"
	}
	if length == "" {
		length = "-1"
	}

	o, err := strconv.Atoi(offset)
	if err != nil {
		r.JSON(400, Message{"message": "offset must be an integer"})
		return
	}

	l, err := strconv.Atoi(length)
	if err != nil {
		r.JSON(400, Message{"message": "length must be an integer"})
		return
	}

	recent := runsList.GetRecent(o, l)
	r.JSON(200, recent)
}

func AddRun(r render.Render, p martini.Params) {
	runsList := models.GetRunList()
	jobsList := models.GetJobList()
	tasksList := models.GetTaskList()

	job, err := jobsList.Get(p["job"])
	if err != nil {
		r.JSON(500, Message{"message": "could not get job"})
		return
	}
	j := job.(models.Job)

	id, err := uuid.NewV4()
	if err != nil {
		r.JSON(500, Message{"message": "could not get new uuid"})
		return
	}

	var tasks []models.Task
	for _, taskName := range j.Tasks {
		task, err := tasksList.Get(taskName)
		if err != nil {
			log.Fatal(err)
			return
		}
		t := task.(models.Task)
		tasks = append(tasks, t)
	}
	err = runsList.AddRun(id.String(), j, tasks)
	if err != nil {
		log.Printf("Could not add run: %v", err)
		r.JSON(500, Message{"message": "could not add run"})
		return
	}

	r.JSON(201, Message{"uuid": id.String()})
}

func GetRun(r render.Render, p martini.Params) {
	runList := models.GetRunList()

	run, err := runList.Get(p["run"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find run"})
		return
	}

	r.JSON(200, run)
}
