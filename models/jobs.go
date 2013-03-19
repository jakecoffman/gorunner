package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Job struct {
	Name  string
	Tasks []Task
}

type JobList struct {
	jobs []Job
}

func (j JobList) GetJobs() []Job {
	return j.jobs
}

func (j JobList) Get(name string) (Job, error) {
	for _, job := range (j.jobs) {
		if job.Name == name {
			return job, nil
		}
	}
	return Job{}, errors.New(fmt.Sprintf("Job '%s' not found", name))
}

func (j *JobList) Append(job Job) error {
	var found bool = false
	for _, j := range(j.jobs) {
		if job.Name == j.Name {
			found = true
		}
	}
	if found {
		return errors.New("Job with that name found in list")
	}
	j.jobs = append(j.jobs, job)
	return nil
}

func (j *JobList) Delete(name string) error {
	var found bool = false
	var i int
	var job Job
	for i, job = range(j.jobs) {
		if job.Name == name {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Job not found for deletion")
	}
	copy(j.jobs[i:], j.jobs[:i])
	j.jobs = j.jobs[:len(j.jobs) - 1]
	return nil
}

func (j JobList) Dumps() string {
	bytes, err := json.Marshal(j.jobs)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (j *JobList) Loads(s string) {
	err := json.Unmarshal([]byte(s), &j.jobs)
	if err != nil {
		panic(err)
	}
}

