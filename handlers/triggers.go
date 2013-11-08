package handlers

import (
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/executor"
	"github.com/jakecoffman/gorunner/models"
	"net/http"
)

func ListTriggers(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(models.Json(triggerList)))
}

func AddTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	payload := unmarshal(r.Body, "name", w)

	trigger := models.Trigger{Name: payload["name"]}
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

	marshal(trigger, w)
}

func UpdateTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	vars := mux.Vars(r)
	trigger, err := models.Get(triggerList, vars["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	payload := unmarshal(r.Body, "cron", w)

	t := trigger.(models.Trigger)
	t.Schedule = payload["cron"]
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
