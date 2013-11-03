package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/executor"
	"github.com/jakecoffman/gorunner/models"
	"github.com/jakecoffman/gorunner/utils"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type addJobPayload struct {
	Name string `json:"name"`
}

func Jobs(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(models.Json(jobList)))
	} else if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var payload addJobPayload
		err = json.Unmarshal(data, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if payload.Name == "" {
			http.Error(w, "Please provide a 'name'", http.StatusBadRequest)
			return
		}

		err = models.Append(jobList, models.Job{Name: payload.Name, Status: "New"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(201)
	} else {
		http.Error(w, fmt.Sprintf("Method '%s' not allowed on this path", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func Job(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := models.Get(jobList, vars["job"])
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	j := job.(models.Job)

	if r.Method == "DELETE" {
		taskPosition, err := strconv.Atoi(vars["task"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		j.DeleteTask(taskPosition)
		models.Update(jobList, j)
	} else if r.Method == "POST" {
		task := r.FormValue("task")
		j.AppendTask(task)
		models.Update(jobList, j)
	}
}

func JobTrigger(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := models.Get(jobList, vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	j := job.(models.Job)

	if r.Method == "DELETE" {
		j.DeleteTrigger(vars["trigger"])
		models.Update(jobList, j)
	} else if r.Method == "POST" {
		trigger := r.FormValue("trigger")
		j.AppendTrigger(trigger)
		triggerList := models.GetTriggerList()
		t, _ := models.Get(triggerList, trigger)
		executor.AddTrigger <- t.(models.Trigger)
		models.Update(jobList, j)
	}

}
