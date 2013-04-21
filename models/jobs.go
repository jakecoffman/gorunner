package models

import (
	"encoding/json"
	"sync"
)

type JobList struct {
	jobs []Job
	sync.RWMutex
}

func (j JobList) GetList() []Job {
	return j.jobs
}

func (j JobList) getList() []Elementer {
	var elements []Elementer
	for _, job := range(j.jobs){
		elements = append(elements, job)
	}
	return elements
}

func (j *JobList) setList(e []Elementer) {
	var jobs []Job
	for _, job := range(e) {
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

