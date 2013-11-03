package models

import (
	"encoding/json"
	"errors"
	"os/exec"
	"sync"
	"time"
)

type Result struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Task   Task      `json:"task"`
	Output string    `json:"output"`
	Error  string    `json:"error"`
}

type Run struct {
	UUID    string    `json:"uuid"`
	Job     Job       `json:"job"`
	Tasks   []Task    `json:"tasks"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Results []Result  `json:"results"`
	Status  string    `json:"status"`
}

func (r Run) ID() string {
	return r.UUID
}

type RunList struct {
	runs []Run
	sync.RWMutex
}

func (j RunList) GetList() []Run {
	return j.runs
}

func (j RunList) getList() []Elementer {
	var elements []Elementer
	for _, run := range j.runs {
		elements = append(elements, run)
	}
	return elements
}

func (j *RunList) setList(e []Elementer) {
	var runs []Run
	for _, run := range e {
		r := run.(Run)
		runs = append(runs, r)
	}
	j.runs = runs
}

func (j *RunList) save() {
	Save(j, runsFile)
}

func (j RunList) Len() int {
	j.RLock()
	defer j.RUnlock()

	return len(j.runs)
}

func (l RunList) Less(i, j int) bool {
	l.RLock()
	defer l.RUnlock()

	return l.runs[i].Start.Before(l.runs[j].Start)
}

func (l RunList) Swap(i, j int) {
	l.RLock()
	defer l.RUnlock()

	l.runs[i], l.runs[j] = l.runs[j], l.runs[i]
}

func (j *RunList) AddRun(UUID string, job Job, tasks []Task) error {
	run := Run{UUID: UUID, Job: job, Tasks: tasks, Start: time.Now(), Status: "New"}
	var found bool = false
	for _, j := range j.runs {
		if run.UUID == j.UUID {
			found = true
		}
	}
	if found {
		return errors.New("Run with that name found in list")
	}
	j.Lock()
	defer j.Unlock()

	j.runs = append(j.runs, run)
	go j.execute(&run)
	Save(&runList, runsFile)
	return nil
}

func (l *RunList) execute(r *Run) {
	r.Status = "Running"
	for _, task := range r.Tasks {
		r.Results = append(r.Results, Result{Start: time.Now(), Task: task})
		result := &r.Results[len(r.Results)-1]
		Update(l, *r)
		cmd := exec.Command("cmd", "/C", task.Script)
		out, err := cmd.Output()
		result.Output = string(out)
		result.End = time.Now()
		if err != nil {
			result.Error = err.Error()
			r.Status = "Failed"
			r.End = time.Now()
			Update(l, *r)
			jobList := GetJobList()
			job, err := Get(jobList, r.Job.Name)
			if err != nil {
				return
			}
			j := job.(Job)
			j.Status = "Failing"
			Update(jobList, job)
			return
		}
		Update(l, *r)
	}
	r.End = time.Now()
	r.Status = "Done"
	jobList := GetJobList()
	job, err := Get(jobList, r.Job.Name)
	if err != nil {
		return
	}
	j := job.(Job)
	j.Status = "Ok"
	Update(jobList, job)
	Update(l, *r)
}

func (j RunList) dumps() string {
	bytes, err := json.Marshal(j.runs)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (j *RunList) loads(s string) {
	err := json.Unmarshal([]byte(s), &j.runs)
	if err != nil {
		panic(err)
	}
}
