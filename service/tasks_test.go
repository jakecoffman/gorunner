package service

import (
	"testing"
)

func TestTaskID(t *testing.T) {
	task := Task{Name: "Task", Script: "echo 'hello'"}
	if task.ID() != "Task" {
		t.Errorf("ID() expected %s but got %s", "Task", task.ID())
	}
}

func TestTaskLoad(t *testing.T) {
	taskList := TaskList{list{elements: make([]elementer, 10), fileName: "test"}}
	value := `[{"name":"mytask","script":"echo 'hi'"}]`
	taskList.Load(mockListReaderFactory(value))
	if string(taskList.dumps()) != value {
		t.Errorf("dumps() expected %s but got %s", value, taskList.dumps())
	}
}
