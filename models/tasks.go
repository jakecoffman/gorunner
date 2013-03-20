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

func (tl *TaskList) Delete(name string) error {
	var i int
	var t Task
	var found bool = false
	for i, t = range(tl.tasks) {
		if t.Name == name {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Could not find task to update")
	}
	tl.tasks = tl.tasks[:i+copy(tl.tasks[i:], tl.tasks[i+1:])]
	return nil
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
