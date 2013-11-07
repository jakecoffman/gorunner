package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/executor"
	"github.com/jakecoffman/gorunner/models"
	"io/ioutil"
	"net/http"
)

type addTriggerPayload struct {
	Name string `json:"name"`
}

type updateTriggerPayload struct {
	Cron string `json:"cron"`
}

func ListTriggers(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(models.Json(triggerList)))
}

func AddTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload addTriggerPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Name == "" {
		http.Error(w, "Please provide a 'name'", http.StatusBadRequest)
		return
	}

	trigger := models.Trigger{Name: payload.Name}
	models.Append(triggerList, trigger)
	w.WriteHeader(201)
}

func GetTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	vars := mux.Vars(r)
	trigger, err := models.Get(triggerList, vars["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(trigger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(bytes)
}

func UpdateTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	vars := mux.Vars(r)
	trigger, err := models.Get(triggerList, vars["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload updateTriggerPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Cron == "" {
		http.Error(w, "Please provide a 'cron'", http.StatusBadRequest)
		return
	}

	t := trigger.(models.Trigger)
	t.Schedule = payload.Cron
	executor.AddTrigger(t)
	err = models.Update(triggerList, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	vars := mux.Vars(r)

	models.Delete(triggerList, vars["trigger"])
}
