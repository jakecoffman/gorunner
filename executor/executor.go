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

func init() {
	AddTrigger = make(chan models.Trigger)
	c = cron.New()
	c.Start()
	c.AddFunc("0 * * * *", func() { fmt.Println("test ran at " + time.Now().Format("2006-01-02 15:04:05")) })
	go func() {
		for {
			select {
			case trigger := <-AddTrigger:
				c.AddFunc(trigger.Schedule, func() { findAndRun(trigger) })
			}
		}
	}()
}

func findAndRun(t models.Trigger) {
	jobList := models.GetJobList()
	for _, job := range jobList.GetList() {
		for _, trigger := range job.Triggers {
			if trigger == t.ID() {
				runnit(job)
				break
			}
		}
	}
}

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
