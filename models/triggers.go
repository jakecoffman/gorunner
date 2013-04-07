package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Trigger struct {
	Name  string
}

type TriggerList struct {
	triggers []Trigger
	lock sync.RWMutex
}

func (j TriggerList) GetList() []Trigger {
	return j.triggers
}

func (j TriggerList) Get(name string) (Trigger, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()

	for _, trigger := range (j.triggers) {
		if trigger.Name == name {
			return trigger, nil
		}
	}
	return Trigger{}, errors.New(fmt.Sprintf("Trigger '%s' not found", name))
}

func (j *TriggerList) Append(name string) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	_, err := j.getPosition(name)
	if err == nil {
		return errors.New("Trigger with that name found in list")
	}
	j.triggers = append(j.triggers, Trigger{name})
	Save(&triggerList, triggersFile)
	return nil
}

func (j *TriggerList) Update(trigger Trigger) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	position, err := j.getPosition(trigger.Name)
	if err != nil {
		return err
	}

	j.triggers[position] = trigger
	Save(&triggerList, triggersFile)
	return nil
}

func (j TriggerList) getPosition(triggerName string) (int,error) {
	var found bool
	var position int
	for i, j := range j.triggers {
		if triggerName == j.Name {
			position = i
			found = true
		}
	}
	if !found {
		return -1, errors.New("Couldn't find " + triggerName)
	}
	return position, nil
}

func (j *TriggerList) Delete(name string) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	var found bool = false
	var i int
	var trigger Trigger
	for i, trigger = range(j.triggers) {
		if trigger.Name == name {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Trigger not found for deletion")
	}
	j.triggers = j.triggers[:i+copy(j.triggers[i:], j.triggers[i+1:])]
	Save(&triggerList, triggersFile)
	return nil
}

func (j TriggerList) Json() string {
	j.lock.RLock()
	defer j.lock.RUnlock()

	return j.dumps()
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

