package models

import (
	"errors"
	"fmt"
)

type Elementer interface {
	ID() string
}

type List interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
	getList() []Elementer
	setList([]Elementer)
	save()
	dumps() string
}

func Get(l List, id string) (Elementer, error) {
	l.RLock()
	defer l.RUnlock()

	list := l.getList()
	for _, job := range (list) {
		if job.ID() == id {
			return job, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Thing '%s' not found", id))
}

func Update(l List, e Elementer) error {
	l.Lock()
	defer l.Unlock()

	position, err := getPosition(l, e.ID())
	if err != nil {
		return err
	}

	things := l.getList()
	things[position] = e
	l.setList(things)
	l.save()
	return nil
}

func Append(l List, e Elementer) error {
	l.Lock()
	defer l.Unlock()

	_, err := getPosition(l, e.ID())
	if err == nil {
		return errors.New("Job with that id found in list")
	}
	l.setList(append(l.getList(), e))
	l.save()
	return nil
}

func Delete(l List, id string) error {
	l.Lock()
	defer l.Unlock()

	list := l.getList()
	var found bool = false
	var i int
	var thing Elementer
	for i, thing = range(list) {
		if thing.ID() == id {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Thing not found for deletion")
	}
	list = list[:i+copy(list[i:], list[i+1:])]
	l.setList(list)
	l.save()
	return nil
}

func Json(l List) string {
	l.RLock()
	defer l.RUnlock()

	return l.dumps()
}

func getPosition(l List, id string) (int,error) {
	var found bool
	var position int
	for i, e := range l.getList() {
		if id == e.ID() {
			position = i
			found = true
		}
	}
	if !found {
		return -1, errors.New("Couldn't find " + id)
	}
	return position, nil
}
