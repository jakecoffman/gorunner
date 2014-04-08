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

func (l *TriggerList) Load(read ListReader) {
	bytes := read(l.fileName)
	var triggers []Trigger
	err := json.Unmarshal([]byte(string(bytes)), &triggers)
	if err != nil {
		panic(err)
	}
	l.elements = []elementer{}
	for _, trigger := range triggers {
		l.elements = append(l.elements, trigger)
	}
}
