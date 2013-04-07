package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/models"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	"github.com/jakecoffman/gorunner/utils"
)

func Jobs(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	if r.Method == "GET" {
		t := template.Must(template.New("_base.html").Funcs(utils.FuncMap).ParseFiles(
			"web/templates/_base.html",
			"web/templates/_nav.html",
			"web/templates/jobs.html",
		))

		if err := t.Execute(w, jobList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		err := jobList.Append(models.Job{Name:name})
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
		t := template.Must(template.New("_base.html").Funcs(utils.FuncMap).ParseFiles(
			"web/templates/_base.html",
			"web/templates/_nav.html",
			"web/templates/job.html",
		))

		if err := t.Execute(w, job); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if r.Method == "POST" {
		task := r.FormValue("task")
		job.Append(task)
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

func JobTask(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	taskPosition, err := strconv.Atoi(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if r.Method == "DELETE" {
		job.Delete(taskPosition)
		jobList.Update(job)
	}
}
