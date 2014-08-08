package service

import (
	"testing"
)

func TestTriggerID(t *testing.T) {
	trigger := Trigger{Name: "Triggy", Schedule: "* * * * * *"}
	if trigger.ID() != "Triggy" {
		t.Errorf("ID() expected %s but got %s", "Triggy", trigger.ID())
	}
}
