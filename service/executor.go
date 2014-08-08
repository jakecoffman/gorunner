package service

import (
	cronService "github.com/jakecoffman/cron"
	"github.com/nu7hatch/gouuid"
)

var triggers map[string]struct{}

type Executor struct {
	cron     *cronService.Cron
	jobList  *JobList
	taskList *TaskList
	runList  *RunList
}

func NewExecutor(jobList *JobList, taskList *TaskList, runList *RunList) *Executor {
	cron := cronService.New()
	cron.Start()
	return &Executor{
		cron,
		jobList,
		taskList,
		runList,
	}
}

func (e Executor) ArmTrigger(t Trigger) {
	e.cron.AddFunc(t.Schedule, func() { e.findAndRun(t) }, t.Name)
}

func (e Executor) DisarmTrigger(name string) {
	e.cron.RemoveJob(name)
	println("Trigger has been removed")
}

// Walks through each job, seeing if the trigger who's turn it is to execute is attached. Executes those jobs.
func (e Executor) findAndRun(t Trigger) {
	jobs := e.jobList.GetJobsWithTrigger(t.ID())
	for _, job := range jobs {
		println("Executing job " + job.Name)
		e.runnit(job)
	}
}

// Gathers the tasks attached to the given job and executes them.
func (e Executor) runnit(j Job) {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	var tasks []Task
	for _, taskName := range j.Tasks {
		task, err := e.taskList.Get(taskName)
		if err != nil {
			panic(err)
		}
		t := task.(Task)
		tasks = append(tasks, t)
	}
	err = e.runList.AddRun(id.String(), j, tasks)
	if err != nil {
		panic(err)
	}
}
