package models

import (
	"testing"
)

func TestID(t *testing.T) {
	trigger := Trigger{Name: "Triggy", Schedule: "* * * * * *"}
	if trigger.ID() != "Triggy" {
		t.Errorf("ID() expected %s but got %s", "Triggy", trigger.ID())
	}
}

func mockListReaderFactory(value string) ListReader {
	return func(string) []byte {
		return []byte(value)
	}
}

func TestLoad(t *testing.T) {
	triggerList = TriggerList{list{elements: make([]elementer, 10), fileName: "somefile.txt"}}
	value := `[{"name":"test","schedule":"* * * * * *"}]`
	triggerList.Load(mockListReaderFactory(value))
	if string(triggerList.dumps()) != value {
		t.Errorf("dumps() expected %s but got %s", value, triggerList.dumps())
	}
}
