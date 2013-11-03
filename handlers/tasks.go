package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/models"
	"io/ioutil"
	"net/http"
)

type addTaskPayload struct {
	Name string `json:"name"`
}

type updateTaskPayload struct {
	Script string `json:"script"`
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(models.Json(taskList)))
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload addTaskPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Name == "" {
		http.Error(w, "Please provide a 'name'", http.StatusBadRequest)
		return
	}

	models.Append(taskList, models.Task{payload.Name, ""})
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

	bytes, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(bytes)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := models.Get(taskList, vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload updateTaskPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Script == "" {
		http.Error(w, "Please provide a 'script'", http.StatusBadRequest)
		return
	}

	t := task.(models.Task)
	t.Script = payload.Script
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
