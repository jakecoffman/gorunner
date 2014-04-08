package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/martini"
	"github.com/jakecoffman/gorunner/executor"
	"github.com/jakecoffman/gorunner/models"
	"github.com/martini-contrib/render"
)

func ListJobs(r render.Render) {
	r.JSON(200, models.GetJobList().List())
}

func AddJob(r render.Render, p martini.Params) {
	jobList := models.GetJobList()

	name := p["name"]

	err := jobList.Append(models.Job{Name: name, Status: "New"})
	if err != nil {
		log.Printf("Error appending to job list: %v", err)
		r.JSON(500, Message{"message": "unable to append to job list"})
		return
	}
	r.JSON(201, Message{"status": "New"})
}

func GetJob(r render.Render, p martini.Params) {
	jobList := models.GetJobList()

	job, err := jobList.Get(p["job"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find job specified"})
		return
	}

	r.JSON(200, job)
}

func DeleteJob(r render.Render, p martini.Params) {
	jobList := models.GetJobList()

	job, err := jobList.Get(p["job"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find job"})
		return
	}

	err = jobList.Delete(job.ID())
	if err != nil {
		r.JSON(500, Message{"message": "could not delete job"})
		return
	}
}

func AddTaskToJob(r render.Render, p martini.Params) {
	jobList := models.GetJobList()

	job, err := jobList.Get(p["job"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find job"})
		return
	}
	j := job.(models.Job)

	j.AppendTask(p["task"])
	jobList.Update(j)

	r.JSON(201, j)
}

func RemoveTaskFromJob(r render.Render, p martini.Params) {
	jobList := models.GetJobList()

	job, err := jobList.Get(p["job"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find job"})
		return
	}
	j := job.(models.Job)

	taskPosition, err := strconv.Atoi(p["task"])
	if err != nil {
		r.JSON(400, Message{"message": "task must be an integer"})
		return
	}
	j.DeleteTask(taskPosition)
	jobList.Update(j)

	r.JSON(200, j)
}

func AddTriggerToJob(r render.Render, p martini.Params) {
	jobList := models.GetJobList()

	job, err := jobList.Get(p["job"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find job"})
		return
	}
	j := job.(models.Job)

	j.AppendTrigger(p["trigger"])
	triggerList := models.GetTriggerList()
	t, err := triggerList.Get(p["trigger"])
	if err != nil {
		r.JSON(500, Message{"message": "unable to add trigger"})
		return
	}
	executor.AddTrigger(t.(models.Trigger))
	jobList.Update(j)

	r.JSON(201, j)
}

func RemoveTriggerFromJob(r render.Render, p martini.Params) {
	jobList := models.GetJobList()

	job, err := jobList.Get(p["job"])
	if err != nil {
		r.JSON(404, Message{"message": "could not find job"})
		return
	}
	j := job.(models.Job)

	j.DeleteTrigger(p["trigger"])
	jobList.Update(j)

	// If Trigger is no longer attached to any Jobs, remove it from Cron to save cycles
	jobs := jobList.GetJobsWithTrigger(p["trigger"])

	if len(jobs) == 0 {
		executor.RemoveTrigger(p["trigger"])
	}

	r.JSON(200, Message{})
}
