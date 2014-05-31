package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/handlers"
	"github.com/jakecoffman/gorunner/hub"
	"github.com/jakecoffman/gorunner/models"
)

const port = "localhost:8090"

var r *mux.Router

// This filter enables messing with the request/response before and after the normal handler
func filter(w http.ResponseWriter, req *http.Request) {
	r.ServeHTTP(w, req) // calls the normal handler
	log.Printf("%s %s %s\n", req.RemoteAddr, req.Method, req.URL)
}

// TODO: Move to handlers package when more websocket handling is required.
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
	r = mux.NewRouter()
	handlers.Install(r)
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
