package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/handlers"
	"net"
	"net/http"
	"os"
)

const port = ":8090"

var r *mux.Router

func app(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/angular/app.html")
}

func gateway(w http.ResponseWriter, req *http.Request) {
	// Before
	r.ServeHTTP(w, req)
	// After
	fmt.Printf("%s %s %s\n", req.RemoteAddr, req.Method, req.URL)
}

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	server := &http.Server{Addr: port, Handler: nil}

	r = mux.NewRouter()

	r.HandleFunc("/", app)
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

	r.HandleFunc("/runs", handlers.ListRuns).Methods("GET")
	r.HandleFunc("/runs", handlers.AddRun).Methods("POST")
	r.HandleFunc("/runs/{run}", handlers.GetRun).Methods("GET")

	r.HandleFunc("/triggers", handlers.ListTriggers).Methods("GET")
	r.HandleFunc("/triggers", handlers.AddTrigger).Methods("POST")
	r.HandleFunc("/triggers/{trigger}", handlers.GetTrigger).Methods("GET")
	r.HandleFunc("/triggers/{trigger}", handlers.UpdateTrigger).Methods("PUT")
	r.HandleFunc("/triggers/{trigger}", handlers.DeleteTrigger).Methods("DELETE")

	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))

	http.HandleFunc("/", gateway)

	go func() {
		for {
			fmt.Println("Running on " + port)
			l, e := net.Listen("tcp", port)
			if e != nil {
				panic(e)
			}
			defer l.Close()
			server.Serve(l)
		}
	}()

	select {}

	println("Dead")
}
