package main

import (
	"github.com/jakecoffman/gorunner/models"
	"github.com/jakecoffman/gorunner/db"
)

func main() {
	var tasks models.TaskList
	t := models.Task{"this", "that"}
	tasks.Append(t)
	println("Wrinting:", tasks.Dumps())
	db.Save(&tasks, "tasks.json")

	// Load and see if equal!
	var newTasks models.TaskList
	db.Load(&newTasks, "tasks.json")
	println("Loaded:", newTasks.Dumps())
}
