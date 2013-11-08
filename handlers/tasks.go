package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/models"
	"io"
	"io/ioutil"
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

func marshal(item interface{}, w http.ResponseWriter) {
	bytes, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(bytes)
}

func unmarshal(r io.Reader, k string, w http.ResponseWriter) (payload map[string]string) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload[k] == "" {
		http.Error(w, "Please provide a '"+k+"'", http.StatusBadRequest)
		return
	}

	return
}
