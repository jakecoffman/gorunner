package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/handlers"
	"github.com/jakecoffman/gorunner/hub"
	"github.com/jakecoffman/gorunner/models"
	"log"
	"net"
	"net/http"
	"os"
)

const port = ":8090"

var r *mux.Router

func setupRoutes() {
	r = mux.NewRouter()

	r.HandleFunc("/", handlers.App)
	r.HandleFunc("/ws", handlers.WsHandler)

	r.HandleFunc("/jobs", handlers.ListJobs).Methods("GET")
	r.HandleFunc("/jobs", handlers.AddJob).Methods("POST")
	r.HandleFunc("/jobs/{job}", handlers.GetJob).Methods("GET")
	r.HandleFunc("/jobs/{job}", handlers.DeleteJob).Methods("DELETE")
	r.HandleFunc("/jobs/{job}/tasks", handlers.AddTaskToJob).Methods("POST")
	r.HandleFunc("/jobs/{job}/tasks/{task}", handlers.RemoveTaskFromJob).Methods("DELETE")
	r.HandleFunc("/jobs/{job}/triggers", handlers.AddTriggerToJob).Methods("POST")
	r.HandleFunc("/jobs/{job}/triggers/{trigger}", handlers.RemoveTriggerFromJob).Methods("DELETE")

	r.HandleFunc("/tasks", handlers.ListTasks).Methods("GET")
	r.HandleFunc("/tasks", handlers.AddTask).Methods("POST")
	r.HandleFunc("/tasks/{task}", handlers.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{task}", handlers.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{task}", handlers.DeleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/{task}/jobs", handlers.ListJobsForTask).Methods("GET")

	r.HandleFunc("/runs", handlers.ListRuns).Methods("GET")
	r.HandleFunc("/runs", handlers.AddRun).Methods("POST")
	r.HandleFunc("/runs/{run}", handlers.GetRun).Methods("GET")

	r.HandleFunc("/triggers", handlers.ListTriggers).Methods("GET")
	r.HandleFunc("/triggers", handlers.AddTrigger).Methods("POST")
	r.HandleFunc("/triggers/{trigger}", handlers.GetTrigger).Methods("GET")
	r.HandleFunc("/triggers/{trigger}", handlers.UpdateTrigger).Methods("PUT")
	r.HandleFunc("/triggers/{trigger}", handlers.DeleteTrigger).Methods("DELETE")
	r.HandleFunc("/triggers/{trigger}/jobs", handlers.ListJobsForTrigger).Methods("GET")

	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))
}

// This filter enables messing with the request/response before and after the normal handler
func filter(w http.ResponseWriter, req *http.Request) {
	r.ServeHTTP(w, req) // calls the normal handler
	log.Printf("%s %s %s\n", req.RemoteAddr, req.Method, req.URL)
}

func getRecentRuns() []byte {
	runsList := models.GetRunListSorted()
	recent := runsList.GetRecent(0, 10)
	bytes, err := json.Marshal(recent)
	if err != nil {
		panic(err.Error())
	}
	return bytes
}

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	hub.NewHub(getRecentRuns)
	go hub.Run()

	// start the server and routes
	server := &http.Server{Addr: port, Handler: nil}
	setupRoutes()
	models.InitDatabase()
	http.HandleFunc("/", filter)

	fmt.Println("Running on " + port)
	l, e := net.Listen("tcp", port)
	if e != nil {
		panic(e)
	}
	defer l.Close()
	server.Serve(l)
}
