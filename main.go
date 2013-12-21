package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/gorunner/handlers"
	"github.com/jakecoffman/gorunner/hub"
	"github.com/jakecoffman/gorunner/models"
	"net"
	"net/http"
	"os"
	"sort"
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := hub.NewConnection(ws)
	hub.Register(c)
	defer hub.Unregister(c)
	go c.Writer()
	c.Reader()
}

func setupRoutes() {
	r = mux.NewRouter()

	r.HandleFunc("/", app)
	r.HandleFunc("/ws", wsHandler)

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

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	hub.NewHub(getRecentRuns)
	go hub.Run()

	// start the server and routes
	server := &http.Server{Addr: port, Handler: nil}
	setupRoutes()
	models.InitDatabase()
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

type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func getRecentRuns() []byte {
	runsList := models.GetRunList()
	sort.Sort(Reverse{runsList})
	recent := runsList.GetRecent(0, 10)
	bytes, err := json.Marshal(recent)
	if err != nil {
		panic(err.Error())
	}
	return bytes
}
