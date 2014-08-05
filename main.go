package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const port = "localhost:8090"

type routeDetail struct {
	route   string
	handler func(http.ResponseWriter, *http.Request)
	method  string
}

var routes []routeDetail = []routeDetail{
	{"/", App, "GET"},
	{"/ws", WsHandler, "GET"},
	{"/jobs", ListJobs, "GET"},
	{"/jobs", AddJob, "POST"},
	{"/jobs/{job}", GetJob, "GET"},
	{"/jobs/{job}", DeleteJob, "DELETE"},
	{"/jobs/{job}/tasks", AddTaskToJob, "POST"},
	{"/jobs/{job}/tasks/{task}", RemoveTaskFromJob, "DELETE"},
	{"/jobs/{job}/triggers/", AddTriggerToJob, "POST"},
	{"/jobs/{job}/triggers/{trigger}", RemoveTriggerFromJob, "DELETE"},

	{"/tasks", ListTasks, "GET"},
	{"/tasks", AddTask, "POST"},
	{"/tasks/{task}", GetTask, "GET"},
	{"/tasks/{task}", UpdateTask, "PUT"},
	{"/tasks/{task}", DeleteTask, "DELETE"},
	{"/tasks/{task}/jobs", ListJobsForTask, "GET"},

	{"/runs", ListRuns, "GET"},
	{"/runs", AddRun, "POST"},
	{"/runs/{run}", GetRun, "GET"},

	{"/triggers", ListTriggers, "GET"},
	{"/triggers", AddTrigger, "POST"},
	{"/triggers/{trigger}", GetTrigger, "GET"},
	{"/triggers/{trigger}", UpdateTrigger, "PUT"},
	{"/triggers/{trigger}", DeleteTrigger, "DELETE"},
	{"/triggers/{trigger}/jobs", ListJobsForTrigger, "GET"},
}

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	NewHub()
	go HubLoop()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))

	for _, detail := range routes {
		r.HandleFunc(detail.route, detail.handler).Methods(detail.method)
	}

	InitDatabase()

	fmt.Println("Running on " + port)
	http.ListenAndServe(port, r)
}
