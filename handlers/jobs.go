package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/db"
	"github.com/jakecoffman/gorunner/models"
	"github.com/gorilla/mux"
	"fmt"
)

const jobsFile = "jobs.json"

func Jobs(w http.ResponseWriter, r *http.Request) {
	var jobList models.JobList
	db.Load(&jobList, jobsFile)

	if r.Method == "GET"{
		;
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		jobList.Append(models.Job{name, []models.Task{}})
		db.Save(&jobList, jobsFile)
	} else {
		http.Error(w, fmt.Sprintf("Method '%s' not allowed on this path", r.Method), http.StatusMethodNotAllowed)
		return
	}

	var jobsTemplate = template.Must(template.ParseFiles(
		"web/templates/_base.html",
		"web/templates/jobs.html",
	))

	if err := jobsTemplate.Execute(w, jobList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Job(w http.ResponseWriter, r *http.Request) {
	var jobList models.JobList
	db.Load(&jobList, jobsFile)
	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		var jobTemplate = template.Must(template.ParseFiles(
			"web/templates/_base.html",
			"web/templates/job.html",
		))

		if err := jobTemplate.Execute(w, job); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if r.Method == "POST" {
		task := r.FormValue("task")
		println(task)
	}


}
