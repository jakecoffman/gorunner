package models

import "encoding/json"

type Job struct {
	Name  string
	Tasks []Task
}

type JobList struct {
	jobs []Job
}

func (j *JobList) Append(job Job){
	j.jobs = append(j.jobs, job)
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

