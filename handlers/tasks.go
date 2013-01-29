package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/db"
	"github.com/jakecoffman/gorunner/models"
)

const tasksFile = "tasks.json"

var tasks = template.Must(template.ParseFiles(
	"web/templates/_base.html",
	"web/templates/tasks.html",
))

func Tasks(w http.ResponseWriter, r *http.Request) {
	var taskList models.TaskList
	db.Load(&taskList, tasksFile)

	if r.Method == "GET"{
		if err := tasks.Execute(w, taskList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		taskList.Append(models.Task{name, ""})
		db.Save(&taskList, tasksFile)
		if err := tasks.Execute(w, taskList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
