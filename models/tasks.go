package models

import (
	"encoding/json"
	"sync"
)

type Task struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

func (t Task) ID() string {
	return t.Name
}

type TaskList struct {
	tasks []Task
	sync.RWMutex
}

func (t TaskList) GetList() []Task {
	return t.tasks
}

func (t TaskList) getList() []Elementer {
	var elements []Elementer
	for _, task := range t.tasks {
		elements = append(elements, task)
	}
	return elements
}

func (t *TaskList) setList(e []Elementer) {
	var tasks []Task
	for _, task := range e {
		t := task.(Task)
		tasks = append(tasks, t)
	}
	t.tasks = tasks
}

func (t *TaskList) save() {
	Save(t, tasksFile)
}

func (t TaskList) dumps() string {
	bytes, err := json.Marshal(t.tasks)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (t *TaskList) loads(s string) {
	err := json.Unmarshal([]byte(s), &t.tasks)
	if err != nil {
		panic(err)
	}
}
