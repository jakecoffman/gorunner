package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	. "github.com/jakecoffman/gorunner/service"
)

const port = "localhost:8090"

var routes = []struct {
	route   string
	handler func(context, http.ResponseWriter, *http.Request) (int, interface{})
	method  string
}{
	{"/jobs", listJobs, "GET"},
	{"/jobs", addJob, "POST"},
	{"/jobs/{job}", getJob, "GET"},
	{"/jobs/{job}", deleteJob, "DELETE"},
	{"/jobs/{job}/tasks", addTaskToJob, "POST"},
	{"/jobs/{job}/tasks/{task}", removeTaskFromJob, "DELETE"},
	{"/jobs/{job}/triggers/", addTriggerToJob, "POST"},
	{"/jobs/{job}/triggers/{trigger}", removeTriggerFromJob, "DELETE"},

	{"/tasks", listTasks, "GET"},
	{"/tasks", addTask, "POST"},
	{"/tasks/{task}", getTask, "GET"},
	{"/tasks/{task}", updateTask, "PUT"},
	{"/tasks/{task}", deleteTask, "DELETE"},
	{"/tasks/{task}/jobs", listJobsForTask, "GET"},

	{"/runs", listRuns, "GET"},
	{"/runs", addRun, "POST"},
	{"/runs/{run}", getRun, "GET"},

	{"/triggers", listTriggers, "GET"},
	{"/triggers", addTrigger, "POST"},
	{"/triggers/{trigger}", getTrigger, "GET"},
	{"/triggers/{trigger}", updateTrigger, "PUT"},
	{"/triggers/{trigger}", deleteTrigger, "DELETE"},
	{"/triggers/{trigger}/jobs", listJobsForTrigger, "GET"},
}

type ctx struct {
	hub         *Hub
	executor    *Executor
	jobList     *JobList
	taskList    *TaskList
	triggerList *TriggerList
	runList     *RunList
}

func (t ctx) Hub() *Hub {
	return t.hub
}

func (t ctx) Executor() *Executor {
	return t.executor
}

func (t ctx) JobList() *JobList {
	return t.jobList
}

func (t ctx) TaskList() *TaskList {
	return t.taskList
}

func (t ctx) TriggerList() *TriggerList {
	return t.triggerList
}

func (t ctx) RunList() *RunList {
	return t.runList
}

type context interface {
	Hub() *Hub
	Executor() *Executor
	JobList() *JobList
	TaskList() *TaskList
	TriggerList() *TriggerList
	RunList() *RunList
}

type appHandler struct {
	*ctx
	handler func(context, http.ResponseWriter, *http.Request) (int, interface{})
}

func (t appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, data := t.handler(t.ctx, w, r)
	marshal(data, w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	log.Println(r.URL, "-", r.Method, "-", code, r.RemoteAddr)
}

func main() {
	wd, _ := os.Getwd()
	log.Println("Working directory", wd)

	jobList := NewJobList()
	taskList := NewTaskList()
	triggerList := NewTriggerList()
	runList := NewRunList(jobList)

	jobList.Load()
	taskList.Load()
	triggerList.Load()
	runList.Load()

	hub := NewHub(runList)
	go hub.HubLoop()

	executor := NewExecutor(jobList, taskList, runList)

	appContext := &ctx{hub, executor, jobList, taskList, triggerList, runList}

	r := mux.NewRouter()

	// non REST routes
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))
	r.HandleFunc("/", app).Methods("GET")
	r.Handle("/ws", appHandler{appContext, wsHandler}).Methods("GET")

	for _, detail := range routes {
		r.Handle(detail.route, appHandler{appContext, detail.handler}).Methods(detail.method)
	}

	log.Println("Running on " + port)
	http.ListenAndServe(port, r)
}
