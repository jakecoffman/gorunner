package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
)

func ListRuns(w http.ResponseWriter, r *http.Request) {
	runsList := GetRunListSorted()

	offset := r.FormValue("offset")
	length := r.FormValue("length")

	if offset == "" {
		offset = "-1"
	}
	if length == "" {
		length = "-1"
	}

	o, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l, err := strconv.Atoi(length)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recent := runsList.GetRecent(o, l)
	marshal(recent, w)
}

func AddRun(w http.ResponseWriter, r *http.Request) {
	runsList := GetRunList()
	jobsList := GetJobList()
	tasksList := GetTaskList()

	payload := unmarshal(r.Body, "job", w)

	job, err := jobsList.Get(payload["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j := job.(Job)

	id, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tasks []Task
	for _, taskName := range j.Tasks {
		task, err := tasksList.Get(taskName)
		if err != nil {
			panic(err)
		}
		t := task.(Task)
		tasks = append(tasks, t)
	}
	err = runsList.AddRun(id.String(), j, tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idResponse := make(map[string]string)
	idResponse["uuid"] = id.String()
	w.WriteHeader(201)
	marshal(idResponse, w)
}

func GetRun(w http.ResponseWriter, r *http.Request) {
	runList := GetRunList()

	vars := mux.Vars(r)
	run, err := runList.Get(vars["run"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(run, w)
}
