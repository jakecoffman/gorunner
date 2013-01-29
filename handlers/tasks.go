package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/db"
	"github.com/jakecoffman/gorunner/models"
	"github.com/gorilla/mux"
)

const tasksFile = "tasks.json"

var taskTemplate = template.Must(template.ParseFiles(
	"web/templates/_base.html",
	"web/templates/task.html",
))
var tasksTemplate = template.Must(template.ParseFiles(
	"web/templates/_base.html",
	"web/templates/tasks.html",
))

func Tasks(w http.ResponseWriter, r *http.Request) {
	var taskList models.TaskList
	db.Load(&taskList, tasksFile)

	if r.Method == "GET"{
		if err := tasksTemplate.Execute(w, taskList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		taskList.Append(models.Task{name, ""})
		db.Save(&taskList, tasksFile)
		if err := tasksTemplate.Execute(w, taskList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Task(w http.ResponseWriter, r *http.Request) {
	var taskList models.TaskList
	db.Load(&taskList, tasksFile)
	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		;
	}

	if err := taskTemplate.Execute(w, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
