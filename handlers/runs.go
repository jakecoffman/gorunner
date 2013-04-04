package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/db"
	"github.com/jakecoffman/gorunner/models"
	"github.com/nu7hatch/gouuid"
	"html/template"
	"net/http"
)

const runsFile = "runs.json"

func Runs(w http.ResponseWriter, r *http.Request) {
	var runsList models.RunList
	db.Load(&runsList, runsFile)

	if r.Method == "GET" {
		var runsTemplate = template.Must(template.ParseFiles(
			"web/templates/_base.html",
			"web/templates/runs.html",
		))

		if err := runsTemplate.Execute(w, runsList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		var jobsList models.JobList
		db.Load(&jobsList, jobsFile)
		var tasksList models.TaskList
		db.Load(&tasksList, tasksFile)

		jobName := r.FormValue("job")
		job, err := jobsList.Get(jobName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		var tasks []models.Task
		for _, taskName := range(job.Tasks){
			task, err := tasksList.Get(taskName)
			if err != nil {
				panic(err)
			}
			tasks = append(tasks, task)
		}
		err = runsList.AddRun(id.String(), job, tasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		db.Save(&runsList, runsFile)
	} else {
		http.Error(w, fmt.Sprintf("Method '%s' not allowed on this path", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func Run(w http.ResponseWriter, r *http.Request) {
	var runList models.RunList
	db.Load(&runList, runsFile)
	vars := mux.Vars(r)
	run, err := runList.Get(vars["run"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		var runTemplate = template.Must(template.ParseFiles(
			"web/templates/_base.html",
			"web/templates/run.html",
		))

		if err := runTemplate.Execute(w, run); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
