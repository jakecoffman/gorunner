package executor

import (
	"fmt"
	"github.com/jakecoffman/gorunner/models"
	"github.com/nu7hatch/gouuid"
	"github.com/robfig/cron"
	"time"
)

var c *cron.Cron
var AddTrigger chan models.Trigger
var triggers map[string]struct{}

func init() {
	triggers = make(map[string]struct{})
	AddTrigger = make(chan models.Trigger)
	c = cron.New()
	c.Start()
	c.AddFunc("0 * * * *", func() { fmt.Println("test ran at " + time.Now().Format("2006-01-02 15:04:05")) })
	go func() {
		for {
			select {
			case trigger := <-AddTrigger:
				_, ok := triggers[trigger.Schedule]
				if !ok {
					triggers[trigger.Schedule] = struct{}{}
					c.AddFunc(trigger.Schedule, func() { findAndRun(trigger) })
				}
			}
		}
	}()
}

// Walks through each job, seeing if the trigger who's turn it is to execute is attached. Executes those jobs.
func findAndRun(t models.Trigger) {
	jobList := models.GetJobList()
	for _, job := range jobList.GetList() {
		for _, trigger := range job.Triggers {
			if trigger == t.ID() {
				fmt.Println("Running job " + job.Name)
				runnit(job)
				break
			}
		}
	}
}

// Gathers the tasks attached to the given job and executes them.
func runnit(j models.Job) {
	tasksList := models.GetTaskList()
	runsList := models.GetRunList()
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	var tasks []models.Task
	for _, taskName := range j.Tasks {
		task, err := models.Get(tasksList, taskName)
		if err != nil {
			panic(err)
		}
		t := task.(models.Task)
		tasks = append(tasks, t)
	}
	err = runsList.AddRun(id.String(), j, tasks)
	if err != nil {
		panic(err)
	}
}
