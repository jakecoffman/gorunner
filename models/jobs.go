package models

import (
	"encoding/json"
	"errors"
	"sync"
)

type Job struct {
	Name     string   `json:"name"`
	Tasks    []string `json:"tasks"`
	Status   string   `json:"status"`
	Triggers []string `json:"triggers"`
}

func (j Job) ID() string {
	return j.Name
}

func (j *Job) AppendTask(task string) {
	j.Tasks = append(j.Tasks, task)
}

func (j *Job) DeleteTask(taskPosition int) error {
	i := taskPosition
	j.Tasks = j.Tasks[:i+copy(j.Tasks[i:], j.Tasks[i+1:])]
	return nil
}

func (j *Job) AppendTrigger(trigger string) error {
	for _, name := range j.Triggers {
		if name == trigger {
			return errors.New("Trigger already on job")
		}
	}
	j.Triggers = append(j.Triggers, trigger)
	return nil
}

func (j *Job) DeleteTrigger(trigger string) error {
	for i, name := range j.Triggers {
		if name == trigger {
			j.Triggers = j.Triggers[:i+copy(j.Triggers[i:], j.Triggers[i+1:])]
			return nil
		}
	}
	return errors.New("Trigger not found")
}

type JobList struct {
	jobs []Job
	sync.RWMutex
}

func (j JobList) GetList() []Job {
	return j.jobs
}

func (j JobList) getList() []Elementer {
	var elements []Elementer
	for _, job := range j.jobs {
		elements = append(elements, job)
	}
	return elements
}

func (j *JobList) setList(e []Elementer) {
	var jobs []Job
	for _, job := range e {
		j := job.(Job)
		jobs = append(jobs, j)
	}
	j.jobs = jobs
}

func (j *JobList) save() {
	Save(j, jobsFile)
}

func (j JobList) dumps() string {
	bytes, err := json.Marshal(j.jobs)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (j *JobList) loads(s string) {
	err := json.Unmarshal([]byte(s), &j.jobs)
	if err != nil {
		panic(err)
	}
}
