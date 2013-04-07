package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"os/exec"
	"sync"
)

type Result struct {
	Start   time.Time
	End     time.Time
	Task    Task
	Output  string
	Error   string
}

type Run struct {
	UUID    string
	Job     Job
	Tasks   []Task
	Start   time.Time
	End     time.Time
	Results []Result
	Status  string
}

type RunList struct {
	runs []Run
	lock sync.RWMutex
}

func (j RunList) GetList() []Run {
	return j.runs
}

func (j RunList) Len() int {
	j.lock.RLock()
	defer j.lock.RUnlock()

	return len(j.runs)
}

func (l RunList) Less(i, j int) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.runs[i].Start.Before(l.runs[j].Start)
}

func (l RunList) Swap(i, j int) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	l.runs[i], l.runs[j] = l.runs[j], l.runs[i]
}

func (j RunList) Get(name string) (*Run, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()

	for _, Run := range (j.runs) {
		if Run.UUID == name {
			return &Run, nil
		}
	}
	return &Run{}, errors.New(fmt.Sprintf("Run '%s' not found", name))
}

func (j *RunList) AddRun(UUID string, job Job, tasks []Task) error {
	run := Run{UUID:UUID, Job:job, Tasks:tasks, Start:time.Now(), Status:"New"}
	var found bool = false
	for _, j := range (j.runs) {
		if run.UUID == j.UUID {
			found = true
		}
	}
	if found {
		return errors.New("Run with that name found in list")
	}
	j.lock.Lock()
	defer j.lock.Unlock()

	j.runs = append(j.runs, run)
	go j.execute(&run)
	Save(&runList, runsFile)
	return nil
}

func (l *RunList) execute(r *Run) {
	r.Status = "Running"
	for _, task := range r.Tasks {
		r.Results = append(r.Results, Result{Start: time.Now(), Task: task})
		result := &r.Results[len(r.Results) - 1]
		l.Update(*r)
		cmd := exec.Command("cmd", "/C", task.Script)
		out, err := cmd.Output()
		result.Output = string(out)
		result.End = time.Now()
		if err != nil {
			result.Error = err.Error()
			r.Status = "Failed"
			r.End = time.Now()
			l.Update(*r)
			return
		}
		l.Update(*r)
	}
	r.End = time.Now()
	r.Status = "Done"
	l.Update(*r)
}

func (j *RunList) Update(run Run) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	var found bool
	var position int
	for i, r := range j.runs {
		if r.UUID == run.UUID {
			position = i
			found = true
		}
	}
	if !found {
		println("Error, can't find run " + run.UUID)
		return errors.New("Can't find run")
	}

	j.runs[position] = run
	Save(&runList, runsFile)
	return nil
}

func (j *RunList) Delete(name string) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	var found bool = false
	var i int
	var Run *Run
	for i, *Run = range (j.runs) {
		if Run.UUID == name {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Run not found for deletion")
	}
	j.runs = j.runs[:i + copy(j.runs[i:], j.runs[i + 1:])]
	Save(&runList, runsFile)
	return nil
}

func (j RunList) Json() string {
	j.lock.RLock()
	defer j.lock.RUnlock()

	return j.dumps()
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

