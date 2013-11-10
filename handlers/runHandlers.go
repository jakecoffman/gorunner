package handlers

import (
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/models"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"sort"
	"strconv"
)

type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func ListRuns(w http.ResponseWriter, r *http.Request) {
	runsList := models.GetRunList()

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

	sort.Sort(Reverse{runsList})
	recent := runsList.GetRecent(o, l)
	marshal(recent, w)
}

func AddRun(w http.ResponseWriter, r *http.Request) {
	runsList := models.GetRunList()
	jobsList := models.GetJobList()
	tasksList := models.GetTaskList()

	payload := unmarshal(r.Body, "job", w)

	job, err := jobsList.Get(payload["job"])
	j := job.(models.Job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tasks []models.Task
	for _, taskName := range j.Tasks {
		task, err := tasksList.Get(taskName)
		if err != nil {
			panic(err)
		}
		t := task.(models.Task)
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
	runList := models.GetRunList()

	vars := mux.Vars(r)
	run, err := runList.Get(vars["run"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(run, w)
}
