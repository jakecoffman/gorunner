package main

import (
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/jakecoffman/gorunner/models"
	"github.com/martini-contrib/render"
)

func ListTasks(r render.Render) {
	r.JSON(200, models.GetTaskList().List())
}

func AddTask(r render.Render, p martini.Params) {
	taskList := models.GetTaskList()

	err := taskList.Append(models.Task{p["name"], ""})
	if err != nil {
		r.JSON(400, Message{"message": "could not append task"})
		return
	}

	r.JSON(201, Message{})
}

func GetTask(r render.Render, p martini.Params) {
	taskList := models.GetTaskList()

	task, err := taskList.Get(p["task"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find task"})
		return
	}

	r.JSON(200, task)
}

func UpdateTask(r render.Render, p martini.Params, req *http.Request) {
	taskList := models.GetTaskList()

	task, err := taskList.Get(p["task"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find task"})
		return
	}

	t := task.(models.Task)
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		r.JSON(400, Message{"message": "bad task payload"})
		return
	}
	t.Script = string(data)
	taskList.Update(t)
	r.JSON(200, t)
}

func DeleteTask(r render.Render, p martini.Params) {
	taskList := models.GetTaskList()

	task, err := taskList.Get(p["task"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find task"})
		return
	}

	taskList.Delete(task.ID())
	r.JSON(200, Message{})
}

func ListJobsForTask(r render.Render, p martini.Params) {
	jobList := models.GetJobList()

	jobs := jobList.GetJobsWithTask(p["task"])

	r.JSON(200, jobs)
}
