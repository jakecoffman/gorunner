package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/handlers"
	"net/http"
)

const port = ":8090"
const base string = "github.com/jakecoffman/gorunner/web/"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.Jobs)
	r.HandleFunc("/jobs", handlers.Jobs)
	r.HandleFunc("/jobs/{job}", handlers.Job)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(base+"static/"))))
	fmt.Println("Running on " + port)
	http.ListenAndServe(port, nil)
}
