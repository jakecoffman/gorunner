package main

import (
	"encoding/json"

	"github.com/codegangsta/martini"
	"github.com/jakecoffman/gorunner/hub"
	"github.com/jakecoffman/gorunner/models"
	"github.com/martini-contrib/render"
)

type Message map[string]interface{}

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
	models.InitDatabase()

	hub.NewHub(getRecentRuns)
	go hub.Run()

	m := martini.Classic()

	m.Use(render.Renderer())

	m.Get("/ws", WsHandler)

	m.Get("/jobs", ListJobs)
	m.Post("/jobs/:name", AddJob)
	m.Get("/jobs/:job", GetJob)
	m.Delete("/jobs/:job", DeleteJob)
	m.Post("/jobs/:job/tasks/:name", AddTaskToJob)
	m.Delete("/jobs/:job/tasks/:task", RemoveTaskFromJob)
	m.Post("/jobs/:job/tasks/:name", AddTriggerToJob)
	m.Delete("/jobs/:job/tasks/:task", RemoveTaskFromJob)

	m.Get("/tasks", ListTasks)
	m.Post("/tasks/:name", AddTask)
	m.Get("/tasks/:task", GetTask)
	m.Put("/tasks/:task", UpdateTask)
	m.Delete("/tasks/:task", DeleteTask)
	m.Get("/tasks/:task/jobs", ListJobsForTask)

	m.Get("/runs", ListRuns)
	m.Post("/runs", AddRun)
	m.Get("/runs/:run", GetRun)

	m.Get("/triggers", ListTriggers)
	m.Post("/triggers/:name", AddTrigger)
	m.Get("/triggers/:trigger", GetTrigger)
	m.Put("/triggers/:trigger", UpdateTrigger)
	m.Delete("/triggers/:trigger", DeleteTrigger)
	m.Get("/triggers/:trigger/jobs", ListJobsForTrigger)

	m.Run()
}
