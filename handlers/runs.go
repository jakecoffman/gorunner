package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/models"
	"github.com/nu7hatch/gouuid"
	"html/template"
	"net/http"
	"sort"
)

type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func Runs(w http.ResponseWriter, r *http.Request) {
	runsList := models.GetRunList()

	if r.Method == "GET" {
		var runsTemplate = template.Must(template.ParseFiles(
			"web/templates/_base.html",
			"web/templates/runs.html",
		))

		sort.Sort(Reverse{runsList})

		if err := runsTemplate.Execute(w, runsList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		jobsList := models.GetJobList()
		tasksList := models.GetTaskList()

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
	} else {
		http.Error(w, fmt.Sprintf("Method '%s' not allowed on this path", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func Run(w http.ResponseWriter, r *http.Request) {
	runList := models.GetRunList()

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
