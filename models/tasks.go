package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Task struct {
	Name   string
	Script string
}

type TaskList struct {
	tasks []Task
	lock sync.RWMutex
}

func (t TaskList) GetList() []Task {
	return t.tasks
}

func (t TaskList) Get(name string) (Task,error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for _, task := range(t.tasks){
		if task.Name == name{
			return task, nil
		}
	}
	return Task{}, errors.New(fmt.Sprintf("Task '%s' not found", name))
}

func (t *TaskList) Append(task Task){
	t.lock.Lock()
	defer t.lock.Unlock()

	t.tasks = append(t.tasks, task)
	Save(&taskList, tasksFile)
}

func (tl *TaskList) Delete(name string) error {
	tl.lock.Lock()
	defer tl.lock.Unlock()

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
	Save(&taskList, tasksFile)
	return nil
}

func (t TaskList) Json() string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.dumps()
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
