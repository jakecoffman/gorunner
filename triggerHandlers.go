package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ListTriggers(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(triggerList.Json()))
}

func AddTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	payload := unmarshal(r.Body, "name", w)

	trigger := Trigger{Name: payload["name"]}
	triggerList.Append(trigger)
	w.WriteHeader(201)
}

func GetTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	vars := mux.Vars(r)
	trigger, err := triggerList.Get(vars["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(trigger, w)
}

func UpdateTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	vars := mux.Vars(r)
	trigger, err := triggerList.Get(vars["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	payload := unmarshal(r.Body, "cron", w)

	t := trigger.(Trigger)
	t.Schedule = payload["cron"]
	ArmTrigger(t)
	err = triggerList.Update(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	vars := mux.Vars(r)

	triggerList.Delete(vars["trigger"])
}

func ListJobsForTrigger(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()
	vars := mux.Vars(r)
	jobs := jobList.GetJobsWithTrigger(vars["trigger"])
	marshal(jobs, w)
}
