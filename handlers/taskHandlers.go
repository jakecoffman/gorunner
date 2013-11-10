package handlers

import (
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/models"
	"net/http"
)

func ListTasks(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(models.Json(taskList)))
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	payload := unmarshal(r.Body, "name", w)

	models.Append(taskList, models.Task{payload["name"], ""})
	w.WriteHeader(201)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := models.Get(taskList, vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(task, w)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := models.Get(taskList, vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	payload := unmarshal(r.Body, "script", w)

	t := task.(models.Task)
	t.Script = payload["script"]
	models.Update(taskList, t)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := models.Get(taskList, vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	models.Delete(taskList, task.ID())
}
