package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/db"
	"github.com/jakecoffman/gorunner/models"
	"github.com/gorilla/mux"
	"strings"
)

const tasksFile = "tasks.json"

func Tasks(w http.ResponseWriter, r *http.Request) {
	var taskList models.TaskList
	db.Load(&taskList, tasksFile)

	var tasksTemplate = template.Must(template.ParseFiles(
		"web/templates/_base.html",
		"web/templates/tasks.html",
	))

	if r.Method == "GET" {
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			if err := tasksTemplate.Execute(w, taskList); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(taskList.Dumps()))
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
		var taskTemplate = template.Must(template.ParseFiles(
			"web/templates/_base.html",
			"web/templates/task.html",
		))

		if err := taskTemplate.Execute(w, task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "PUT" {
		script := r.FormValue("script")
		task.Script = script
		taskList.Delete(task.Name)
		taskList.Append(task)
		db.Save(&taskList, tasksFile)
	} else if r.Method == "DELETE" {
		taskList.Delete(task.Name)
		db.Save(&taskList, tasksFile)
	} else {
		http.Error(w, "Unknown method type" , http.StatusMethodNotAllowed)
	}


}
