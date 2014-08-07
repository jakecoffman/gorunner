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
	r.HandleFunc("/", app).Methods("GET")
	r.Handle("/ws", appHandler{appContext, wsHandler}).Methods("GET")

	for _, detail := range routes {
		r.Handle(detail.route, appHandler{appContext, detail.handler}).Methods(detail.method)
	}

	log.Println("Running on " + port)
	http.ListenAndServe(port, r)
}
