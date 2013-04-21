package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/models"
	"github.com/jakecoffman/gorunner/utils"
	"github.com/nu7hatch/gouuid"
	"html/template"
	"net/http"
	"sort"
	"strings"
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
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			t := template.Must(template.New("_base.html").Funcs(utils.FuncMap).ParseFiles(
				"web/templates/_base.html",
				"web/templates/_nav.html",
				"web/templates/runs.html",
			))

			sort.Sort(Reverse{runsList})

			if err := t.Execute(w, runsList); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			sort.Sort(Reverse{runsList})
			w.Write([]byte(models.Json(runsList)))
		}
	} else if r.Method == "POST" {
		jobsList := models.GetJobList()
		tasksList := models.GetTaskList()

		jobName := r.FormValue("job")
		job, err := models.Get(jobsList, jobName)
		j := job.(models.Job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		var tasks []models.Task
		for _, taskName := range(j.Tasks){
			task, err := models.Get(tasksList, taskName)
			if err != nil {
				panic(err)
			}
			t := task.(models.Task)
			tasks = append(tasks, t)
		}
		err = runsList.AddRun(id.String(), j, tasks)
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
	run, err := models.Get(runList, vars["run"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		t := template.Must(template.New("_base.html").Funcs(utils.FuncMap).ParseFiles(
			"web/templates/_base.html",
			"web/templates/_nav.html",
			"web/templates/run.html",
		))

		if err := t.Execute(w, run); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
