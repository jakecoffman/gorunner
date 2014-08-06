package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const port = "localhost:8090"

var routes = []struct {
	route   string
	handler func(*context, http.ResponseWriter, *http.Request)
	method  string
}{
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

type context struct {
	hub         *Hub
	executor    *Executor
	jobList     *JobList
	taskList    *TaskList
	triggerList *TriggerList
	runList     *RunList
}

type appHandler struct {
	*context
	handler func(*context, http.ResponseWriter, *http.Request)
}

func (t appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.handler(t.context, w, r)
	log.Println(r.URL, r.Method, r.RemoteAddr)
}

func main() {
	wd, _ := os.Getwd()
	log.Println("Working directory", wd)

	jobList := &JobList{
		list{elements: []elementer{}, fileName: jobsFile},
	}
	taskList := &TaskList{
		list{elements: []elementer{}, fileName: tasksFile},
	}
	triggerList := &TriggerList{
		list{elements: []elementer{}, fileName: triggersFile},
	}
	runList := &RunList{
		list{elements: []elementer{}, fileName: runsFile},
		jobList,
	}

	jobList.Load()
	taskList.Load(readFile)
	triggerList.Load(readFile)
	runList.Load()

	hub := NewHub(runList)
	go hub.HubLoop()

	executor := NewExecutor(jobList, taskList, runList)

	appContext := &context{hub, executor, jobList, taskList, triggerList, runList}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))

	for _, detail := range routes {
		r.Handle(detail.route, appHandler{appContext, detail.handler}).Methods(detail.method)
	}

	log.Println("Running on " + port)
	http.ListenAndServe(port, r)
}
