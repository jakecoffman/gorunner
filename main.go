package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/handlers"
	"net/http"
	"os"
)

const port = ":8090"

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.Jobs)
	r.HandleFunc("/jobs", handlers.Jobs)
	r.HandleFunc("/jobs/{job}", handlers.Job)

	r.HandleFunc("/tasks", handlers.Tasks)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	fmt.Println("Running on " + port)
	http.ListenAndServe(port, nil)
}
