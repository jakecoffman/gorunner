package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Job struct {
	Name  string
	Tasks []string
}

func (j *Job) Delete(taskPosition int) error {
	i := taskPosition
	j.Tasks = j.Tasks[:i+copy(j.Tasks[i:], j.Tasks[i+1:])]
	return nil
}

type JobList struct {
	jobs []Job
	lock sync.RWMutex
}

func GetJobList() *JobList {
	return &jobList
}

func (j JobList) GetList() []Job {
	return j.jobs
}

func (j JobList) Get(name string) (Job, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()

	for _, job := range (j.jobs) {
		if job.Name == name {
			return job, nil
		}
	}
	return Job{}, errors.New(fmt.Sprintf("Job '%s' not found", name))
}

func (j *JobList) Append(job Job) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	_, err := j.getPosition(job.Name)
	if err == nil {
		return errors.New("Job with that name found in list")
	}
	j.jobs = append(j.jobs, job)
	Save(&jobList, jobsFile)
	return nil
}

func (j *JobList) Update(job Job) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	position, err := j.getPosition(job.Name)
	if err != nil {
		return err
	}

	j.jobs[position] = job
	Save(&jobList, jobsFile)
	return nil
}

func (j JobList) getPosition(jobName string) (int,error) {
	var found bool
	var position int
	for i, j := range j.jobs {
		if jobName == j.Name {
			position = i
			found = true
		}
	}
	if !found {
		return -1, errors.New("Couldn't find " + jobName)
	}
	return position, nil
}

func (j *JobList) Delete(name string) error {
	j.lock.Lock()
	defer j.lock.Unlock()

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
	j.jobs = j.jobs[:i+copy(j.jobs[i:], j.jobs[i+1:])]
	Save(&jobList, jobsFile)
	return nil
}

func (j JobList) Json() string {
	j.lock.RLock()
	defer j.lock.RUnlock()

	return j.dumps()
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

