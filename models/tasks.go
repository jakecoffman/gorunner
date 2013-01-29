package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Task struct {
	Name   string
	Script string
}

type TaskList struct {
	tasks []Task
}

func (t TaskList) GetTasks() []Task{
	return t.tasks
}

func (t TaskList) Get(name string) (Task,error) {
	for _, task := range(t.tasks){
		if task.Name == name{
			return task, nil
		}
	}
	return Task{}, errors.New(fmt.Sprintf("Task '%s' not found", name))
}

func (t *TaskList) Append(task Task){
	t.tasks = append(t.tasks, task)
}

func (t TaskList) Dumps() string {
	bytes, err := json.Marshal(t.tasks)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (t *TaskList) Loads(s string) {
	err := json.Unmarshal([]byte(s), &t.tasks)
	if err != nil {
		panic(err)
	}
}
