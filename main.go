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
	handler func(*context, http.ResponseWriter, *http.Request) (int, interface{})
	method  string
}{
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
	handler func(*context, http.ResponseWriter, *http.Request) (int, interface{})
}

func (t appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, data := t.handler(t.context, w, r)
	marshal(data, w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	log.Println(r.URL, "-", r.Method, "-", code, r.RemoteAddr)
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

	// non REST routes
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))
	r.HandleFunc("/", App).Methods("GET")
	r.Handle("/ws", appHandler{appContext, WsHandler}).Methods("GET")

	for _, detail := range routes {
		r.Handle(detail.route, appHandler{appContext, detail.handler}).Methods(detail.method)
	}

	log.Println("Running on " + port)
	http.ListenAndServe(port, r)
}
