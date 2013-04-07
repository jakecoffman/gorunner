package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/handlers"
	"github.com/jakecoffman/gorunner/executor"
	"net/http"
	"os"
	"net"
)

const port = ":8090"

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	server := &http.Server{Addr: port, Handler: nil }
	l, e := net.Listen("tcp", port)
	if e != nil {
		panic(e)
	}

	r := mux.NewRouter()
	r = mux.NewRouter()

	r.HandleFunc("/", handlers.Jobs)
	r.HandleFunc("/jobs", handlers.Jobs)
	r.HandleFunc("/jobs/{job}", handlers.Job)
	r.HandleFunc("/jobs/{job}/{task}", handlers.JobTask)

	r.HandleFunc("/tasks", handlers.Tasks)
	r.HandleFunc("/tasks/{task}", handlers.Task)

	r.HandleFunc("/runs", handlers.Runs)
	r.HandleFunc("/runs/{run}", handlers.Run)

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

	<-executor.Kill
	println("Dead")
}
