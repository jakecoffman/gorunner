package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/models"
	"github.com/gorilla/mux"
	"fmt"
)

const jobsFile = "jobs.json"

func Jobs(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	if r.Method == "GET" {
		var jobsTemplate = template.Must(template.ParseFiles(
			"web/templates/_base.html",
			"web/templates/jobs.html",
		))

		if err := jobsTemplate.Execute(w, jobList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		err := jobList.Append(models.Job{name, []string{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, fmt.Sprintf("Method '%s' not allowed on this path", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func Job(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

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
		job.Tasks = append(job.Tasks, task)
		jobList.Delete(job.Name)
		jobList.Append(job)
	} else if r.Method == "DELETE" {
		err := jobList.Delete(job.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}


}
