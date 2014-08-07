package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type elementer interface {
	ID() string
}

type list struct {
	elements []elementer
	fileName string
	sync.RWMutex
}

func (l *list) Get(id string) (elementer, error) {
	l.RLock()
	defer l.RUnlock()

	for _, e := range l.elements {
		if e.ID() == id {
			return e, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Element '%s' not found", id))
}

func (l *list) Update(e elementer) error {
	l.Lock()
	defer l.Unlock()

	position, err := l.pos(e.ID())
	if err != nil {
		return err
	}

	l.elements[position] = e
	l.save()
	return nil
}

func (l *list) Append(e elementer) error {
	l.Lock()
	defer l.Unlock()

	if e.ID() == "" {
		return errors.New("No ID provided")
	}

	_, err := l.pos(e.ID())
	if err == nil {
		return errors.New("Element with that id found in list")
	}
	l.elements = append(l.elements, e)
	l.save()
	return nil
}

func (l *list) Delete(id string) error {
	l.Lock()
	defer l.Unlock()

	var found bool = false
	var i int
	var thing elementer
	for i, thing = range l.elements {
		if thing.ID() == id {
			found = true
			break
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("Element '%s' not found for deletion", id))
	}
	l.elements = l.elements[:i+copy(l.elements[i:], l.elements[i+1:])]
	l.save()
	return nil
}

func (l *list) Json() []byte {
	l.RLock()
	defer l.RUnlock()

	return l.dumps()
}

func (l list) save() {
	writeFile(l.dumps(), l.fileName)
}

func (l list) dumps() []byte {
	bytes, err := json.Marshal(l.elements)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (l list) Dump() []elementer {
	l.RLock()
	defer l.RUnlock()
	return l.elements
}

func (l list) pos(id string) (int, error) {
	for i, e := range l.elements {
		if e.ID() == id {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}
