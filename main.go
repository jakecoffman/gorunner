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

func app(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/angular/app.html")
}

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	server := &http.Server{Addr: port, Handler: nil}

	r := mux.NewRouter()
	r = mux.NewRouter()

	r.HandleFunc("/", app)
	r.HandleFunc("/jobs", handlers.ListJobs).Methods("GET")
	r.HandleFunc("/jobs", handlers.AddJob).Methods("POST")
	r.HandleFunc("/jobs/{job}", handlers.GetJob).Methods("GET")
	r.HandleFunc("/jobs/{job}", handlers.DeleteJob).Methods("DELETE")
	r.HandleFunc("/jobs/{job}/tasks", handlers.JobTask)
	r.HandleFunc("/jobs/{job}/tasks/{task}", handlers.JobTask)
	r.HandleFunc("/jobs/{job}/triggers", handlers.JobTrigger)
	r.HandleFunc("/jobs/{job}/triggers/{trigger}", handlers.JobTrigger)

	r.HandleFunc("/tasks", handlers.ListTasks).Methods("GET")
	r.HandleFunc("/tasks", handlers.AddTask).Methods("POST")
	r.HandleFunc("/tasks/{task}", handlers.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{task}", handlers.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{task}", handlers.DeleteTask).Methods("DELETE")

	r.HandleFunc("/runs", handlers.Runs)
	r.HandleFunc("/runs/{run}", handlers.Run)

	r.HandleFunc("/triggers", handlers.Triggers)
	r.HandleFunc("/triggers/{trigger}", handlers.Trigger)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

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
