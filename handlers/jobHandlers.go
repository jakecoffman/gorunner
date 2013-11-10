package handlers

import (
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/executor"
	"github.com/jakecoffman/gorunner/models"
	"net/http"
	"strconv"
)

func ListJobs(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	w.Write([]byte(jobList.Json()))
}

func AddJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	payload := unmarshal(r.Body, "name", w)

	err := jobList.Append(models.Job{Name: payload["name"], Status: "New"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(job, w)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = jobList.Delete(job.ID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddTaskToJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(models.Job)

	payload := unmarshal(r.Body, "task", w)
	j.AppendTask(payload["task"])
	jobList.Update(j)

	w.WriteHeader(201)
}

func RemoveTaskFromJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(models.Job)

	taskPosition, err := strconv.Atoi(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j.DeleteTask(taskPosition)
	jobList.Update(j)
}

func AddTriggerToJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(models.Job)

	payload := unmarshal(r.Body, "trigger", w)

	j.AppendTrigger(payload["trigger"])
	triggerList := models.GetTriggerList()
	t, err := triggerList.Get(payload["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	executor.AddTrigger(t.(models.Trigger))
	jobList.Update(j)

	w.WriteHeader(201)
}

func RemoveTriggerFromJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(models.Job)

	t := vars["trigger"]
	j.DeleteTrigger(t)

	// If Trigger is no longer attached to any Jobs, remove it from Cron
	found := false
	for _, e := range jobList.GetAll() {
		job := e.(models.Job)
		// Exclude this job because we just removed it. Kind of a race.
		if job.Name == j.Name {
			continue
		}

		for _, trigger := range job.Triggers {
			if trigger == t {
				found = true
				break
			}
		}
	}

	if found == false {
		executor.RemoveTrigger(t)
	}

	jobList.Update(j)
}
