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
