package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/models"
	"github.com/jakecoffman/gorunner/utils"
	"github.com/gorilla/mux"
	"strings"
)

func Tasks(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	if r.Method == "GET" {
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			t := template.Must(template.New("_base.html").Funcs(utils.FuncMap).ParseFiles(
				"web/templates/_base.html",
				"web/templates/_nav.html",
				"web/templates/tasks.html",
			))

			if err := t.Execute(w, taskList); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(models.Json(taskList)))
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		models.Append(taskList, models.Task{name, ""})
	} else {
		http.Error(w, "Unknown method type" , http.StatusMethodNotAllowed)
	}
}

func Task(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := models.Get(taskList, vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		t := template.Must(template.New("_base.html").Funcs(utils.FuncMap).ParseFiles(
			"web/templates/_base.html",
			"web/templates/_nav.html",
			"web/templates/task.html",
		))

		if err := t.Execute(w, task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "PUT" {
		script := r.FormValue("script")
		t := task.(models.Task)
		t.Script = script
		models.Update(taskList, t)
	} else if r.Method == "DELETE" {
		models.Delete(taskList, task.ID())
	} else {
		http.Error(w, "Unknown method type" , http.StatusMethodNotAllowed)
	}
}
