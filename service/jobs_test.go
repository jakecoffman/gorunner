package service

import (
	"fmt"
	"testing"
)

func TestJobID(t *testing.T) {
	job := Job{Name: "name"}
	if job.ID() != "name" {
		t.Errorf("ID() expected %s but got %s", "name", job.ID())
	}
}

func TestJobAppendTask(t *testing.T) {
	job := Job{"name", make([]string, 0), "status", make([]string, 0)}
	job.AppendTask("task")
	expected := []string{"task"}
	if fmt.Sprintf("%#v", job.Tasks) != fmt.Sprintf("%#v", expected) {
		t.Errorf("Expected %#v but got %#v", expected, job.Tasks)
	}
}
