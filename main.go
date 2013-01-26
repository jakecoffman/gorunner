package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/handlers"
	"net/http"
)

const Port = ":8090"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.Jobs)
	r.HandleFunc("/jobs", handlers.Jobs)
	r.HandleFunc("/jobs/{job}", handlers.Job)

	http.Handle("/", r)
	fmt.Println("Running on " + Port)
	http.ListenAndServe(Port, nil)
}
