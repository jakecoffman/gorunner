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
	l, e := net.Listen("tcp", port)
	if e != nil {
		panic(e)
	}

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

	r.HandleFunc("/tasks", handlers.Tasks)
	r.HandleFunc("/tasks/{task}", handlers.Task)

	r.HandleFunc("/runs", handlers.Runs)
	r.HandleFunc("/runs/{run}", handlers.Run)

	r.HandleFunc("/triggers", handlers.Triggers)
	r.HandleFunc("/triggers/{trigger}", handlers.Trigger)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	go func() {
		// run stuff outside of the main loop
	}()

	go func() {
		for {
			fmt.Println("Running on " + port)
			server.Serve(l)
			l, e = net.Listen("tcp", port)
			if e != nil {
				panic(e)
			}
		}
	}()
	defer l.Close()

	select {}

	println("Dead")
}
