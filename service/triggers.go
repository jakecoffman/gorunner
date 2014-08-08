package service

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

func NewTriggerList() *TriggerList {
	return &TriggerList{
		list{elements: []elementer{}, fileName: triggersFile},
	}
}

func (l *TriggerList) Load() {
	bytes := readFile(l.fileName)
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
