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
		err := models.Append(jobList, models.Job{Name:name, Status:"New"})
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
	job, err := models.Get(jobList, vars["job"])
	j := job.(models.Job)
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
		j.Append(task)
		models.Update(jobList, j)
	} else if r.Method == "DELETE" {
		err := models.Delete(jobList, job.ID())
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
	job, err := models.Get(jobList, vars["job"])
	j := job.(models.Job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	taskPosition, err := strconv.Atoi(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if r.Method == "DELETE" {
		j.Delete(taskPosition)
		models.Update(jobList, j)
	}
}
