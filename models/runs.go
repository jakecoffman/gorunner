package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Run struct {
	UUID  string
	Job   Job
	Tasks []Task
	Start time.Time
	End   time.Time
}

type RunList struct {
	Runs []Run
}

func (j RunList) GetRuns() []Run {
	return j.Runs
}

func (j RunList) Get(name string) (Run, error) {
	for _, Run := range (j.Runs) {
		if Run.UUID == name {
			return Run, nil
		}
	}
	return Run{}, errors.New(fmt.Sprintf("Run '%s' not found", name))
}

func (j *RunList) Append(run Run) error {
	var found bool = false
	for _, j := range (j.Runs) {
		if run.UUID == j.UUID {
			found = true
		}
	}
	if found {
		return errors.New("Run with that name found in list")
	}
	run.Start = time.Now()
	j.Runs = append(j.Runs, run)
	return nil
}

func (j *RunList) Delete(name string) error {
	var found bool = false
	var i int
	var Run Run
	for i, Run = range (j.Runs) {
		if Run.UUID == name {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Run not found for deletion")
	}
	j.Runs = j.Runs[:i + copy(j.Runs[i:], j.Runs[i + 1:])]
	return nil
}

func (j RunList) Dumps() string {
	bytes, err := json.Marshal(j.Runs)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (j *RunList) Loads(s string) {
	err := json.Unmarshal([]byte(s), &j.Runs)
	if err != nil {
		panic(err)
	}
}

