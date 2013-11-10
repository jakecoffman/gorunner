package models

import (
	"encoding/json"
)

type Task struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

func (t Task) ID() string {
	return t.Name
}

type TaskList struct {
	list
}

func (l TaskList) Save() {
	var tasks []Task

	for _, e := range l.elements {
		tasks = append(tasks, e.(Task))
	}

	bytes, err := json.Marshal(tasks)
	if err != nil {
		panic(err)
	}
	writeFile(bytes, l.fileName)
}

func (l *TaskList) Load() {
	bytes := readFile(l.fileName)
	var tasks []Task
	err := json.Unmarshal([]byte(string(bytes)), &tasks)
	if err != nil {
		panic(err)
	}
	l.elements = nil
	for _, task := range tasks {
		l.elements = append(l.elements, task)
	}
}
