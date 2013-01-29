package models

import "encoding/json"

type Task struct {
	Name   string
	Script string
}

type TaskList struct {
	tasks []Task
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
