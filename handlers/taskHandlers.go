package handlers

import (
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/models"
	"net/http"
)

func ListTasks(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(taskList.Json()))
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	payload := unmarshal(r.Body, "name", w)

	err := taskList.Append(models.Task{payload["name"], ""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(201)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(task, w)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	payload := unmarshal(r.Body, "script", w)

	t := task.(models.Task)
	t.Script = payload["script"]
	taskList.Update(t)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	taskList.Delete(task.ID())
}

func ListJobsForTask(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()
	vars := mux.Vars(r)
	jobs := jobList.GetJobsWithTask(vars["task"])
	marshal(jobs, w)
}
