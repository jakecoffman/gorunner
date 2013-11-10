package models

import (
	"encoding/json"
)

type Trigger struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
}

func (t Trigger) ID() string {
	return t.Name
}

type TriggerList struct {
	list
}

func (l TriggerList) Save() {
	var triggers []Trigger

	for _, e := range l.elements {
		triggers = append(triggers, e.(Trigger))
	}

	bytes, err := json.Marshal(triggers)
	if err != nil {
		panic(err)
	}
	writeFile(bytes, l.fileName)
}

func (l *TriggerList) Load() {
	bytes := readFile(l.fileName)
	var triggers []Trigger
	err := json.Unmarshal([]byte(string(bytes)), &triggers)
	if err != nil {
		panic(err)
	}
	l.elements = nil
	for _, trigger := range triggers {
		l.elements = append(l.elements, trigger)
	}
}
