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

var jobsTemplate = template.Must(template.ParseFiles(
	"web/templates/_base.html",
	"web/templates/jobs.html",
))

var jobTemplate = template.Must(template.ParseFiles(
	"web/templates/_base.html",
	"web/templates/job.html",
))

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
		;
	}

	if err := jobTemplate.Execute(w, job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
