package models

import (
	"encoding/json"
	"sync"
)

type Trigger struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
}

func (t Trigger) ID() string {
	return t.Name
}

type TriggerList struct {
	triggers []Trigger
	sync.RWMutex
}

func (t TriggerList) GetList() []Trigger {
	return t.triggers
}

func (t TriggerList) getList() []Elementer {
	var elements []Elementer
	for _, trigger := range t.triggers {
		elements = append(elements, trigger)
	}
	return elements
}

func (t *TriggerList) setList(e []Elementer) {
	var triggers []Trigger
	for _, trigger := range e {
		t := trigger.(Trigger)
		triggers = append(triggers, t)
	}
	t.triggers = triggers
}

func (t *TriggerList) save() {
	Save(t, triggersFile)
}

func (j TriggerList) dumps() string {
	bytes, err := json.Marshal(j.triggers)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (j *TriggerList) loads(s string) {
	err := json.Unmarshal([]byte(s), &j.triggers)
	if err != nil {
		panic(err)
	}
}
