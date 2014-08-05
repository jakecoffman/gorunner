package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

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

func Install(r *mux.Router) {
	r.HandleFunc("/", App)
	r.HandleFunc("/ws", WsHandler)

	r.HandleFunc("/jobs", ListJobs).Methods("GET")
	r.HandleFunc("/jobs", AddJob).Methods("POST")
	r.HandleFunc("/jobs/{job}", GetJob).Methods("GET")
	r.HandleFunc("/jobs/{job}", DeleteJob).Methods("DELETE")
	r.HandleFunc("/jobs/{job}/tasks", AddTaskToJob).Methods("POST")
	r.HandleFunc("/jobs/{job}/tasks/{task}", RemoveTaskFromJob).Methods("DELETE")
	r.HandleFunc("/jobs/{job}/triggers", AddTriggerToJob).Methods("POST")
	r.HandleFunc("/jobs/{job}/triggers/{trigger}", RemoveTriggerFromJob).Methods("DELETE")

	r.HandleFunc("/tasks", ListTasks).Methods("GET")
	r.HandleFunc("/tasks", AddTask).Methods("POST")
	r.HandleFunc("/tasks/{task}", GetTask).Methods("GET")
	r.HandleFunc("/tasks/{task}", UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{task}", DeleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/{task}/jobs", ListJobsForTask).Methods("GET")

	r.HandleFunc("/runs", ListRuns).Methods("GET")
	r.HandleFunc("/runs", AddRun).Methods("POST")
	r.HandleFunc("/runs/{run}", GetRun).Methods("GET")

	r.HandleFunc("/triggers", ListTriggers).Methods("GET")
	r.HandleFunc("/triggers", AddTrigger).Methods("POST")
	r.HandleFunc("/triggers/{trigger}", GetTrigger).Methods("GET")
	r.HandleFunc("/triggers/{trigger}", UpdateTrigger).Methods("PUT")
	r.HandleFunc("/triggers/{trigger}", DeleteTrigger).Methods("DELETE")
	r.HandleFunc("/triggers/{trigger}/jobs", ListJobsForTrigger).Methods("GET")

	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))
}
