package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/db"
	"github.com/jakecoffman/gorunner/models"
	"github.com/jakecoffman/gorunner/execution"
	"github.com/gorilla/mux"
)

const tasksFile = "tasks.json"

func Tasks(w http.ResponseWriter, r *http.Request) {
	var taskList models.TaskList
	db.Load(&taskList, tasksFile)

	var tasksTemplate = template.Must(template.ParseFiles(
		"web/templates/_base.html",
		"web/templates/tasks.html",
	))

	execution.Kill<- true
	return

	if r.Method == "GET" {
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
	} else {
		http.Error(w, "Unknown method type" , http.StatusMethodNotAllowed)
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
	} else {
		http.Error(w, "Unknown method type" , http.StatusMethodNotAllowed)
	}

	var taskTemplate = template.Must(template.ParseFiles(
		"web/templates/_base.html",
		"web/templates/task.html",
	))

	if err := taskTemplate.Execute(w, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
