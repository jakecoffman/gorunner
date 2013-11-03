package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/gorunner/models"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"net/http"
	"sort"
)

type Reverse struct {
	sort.Interface
}

type addRunPayload struct {
	Job string `json:"job"`
}

type addRunResponse struct {
	Uuid string `json:"uuid"`
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func ListRuns(w http.ResponseWriter, r *http.Request) {
	runsList := models.GetRunList()

	sort.Sort(Reverse{runsList})
	w.Write([]byte(models.Json(runsList)))
}

func AddRun(w http.ResponseWriter, r *http.Request) {
	runsList := models.GetRunList()
	jobsList := models.GetJobList()
	tasksList := models.GetTaskList()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload addRunPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Job == "" {
		http.Error(w, "Please provide a 'job' to run", http.StatusBadRequest)
		return
	}

	job, err := models.Get(jobsList, payload.Job)
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
		task, err := models.Get(tasksList, taskName)
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

	var idResponse addRunResponse
	idResponse.Uuid = id.String()
	data, err = json.Marshal(idResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	w.Write(data)
}

func GetRun(w http.ResponseWriter, r *http.Request) {
	runList := models.GetRunList()

	vars := mux.Vars(r)
	run, err := models.Get(runList, vars["run"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(run)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(bytes)
}
