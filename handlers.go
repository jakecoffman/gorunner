package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
)

// General

func App(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/app.html")
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := NewConnection(ws)
	Register(c)
	defer Unregister(c)
	go c.Writer()
	c.Reader()
}

// Jobs

func ListJobs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(GetJobList().Json()))
}

func AddJob(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()

	payload := unmarshal(r.Body, "name", w)

	err := jobList.Append(Job{Name: payload["name"], Status: "New"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(job, w)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = jobList.Delete(job.ID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddTaskToJob(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(Job)

	payload := unmarshal(r.Body, "task", w)
	j.AppendTask(payload["task"])
	jobList.Update(j)

	w.WriteHeader(201)
}

func RemoveTaskFromJob(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(Job)

	taskPosition, err := strconv.Atoi(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j.DeleteTask(taskPosition)
	jobList.Update(j)
}

func AddTriggerToJob(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(Job)

	payload := unmarshal(r.Body, "trigger", w)

	j.AppendTrigger(payload["trigger"])
	triggerList := GetTriggerList()
	t, err := triggerList.Get(payload["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ArmTrigger(t.(Trigger))
	jobList.Update(j)

	w.WriteHeader(201)
}

func RemoveTriggerFromJob(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(Job)

	t := vars["trigger"]
	j.DeleteTrigger(t)
	jobList.Update(j)

	// If Trigger is no longer attached to any Jobs, remove it from Cron to save cycles
	jobs := jobList.GetJobsWithTrigger(t)

	if len(jobs) == 0 {
		DisarmTrigger(t)
	}
}

// Run

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

// Tasks
func ListTasks(w http.ResponseWriter, r *http.Request) {
	taskList := GetTaskList()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(taskList.Json()))
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	taskList := GetTaskList()

	payload := unmarshal(r.Body, "name", w)

	err := taskList.Append(Task{payload["name"], ""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(201)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	taskList := GetTaskList()

	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(task, w)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskList := GetTaskList()

	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	payload := unmarshal(r.Body, "script", w)

	t := task.(Task)
	t.Script = payload["script"]
	taskList.Update(t)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskList := GetTaskList()

	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	taskList.Delete(task.ID())
}

func ListJobsForTask(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()
	vars := mux.Vars(r)
	jobs := jobList.GetJobsWithTask(vars["task"])
	marshal(jobs, w)
}

// Triggers

func ListTriggers(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(triggerList.Json()))
}

func AddTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	payload := unmarshal(r.Body, "name", w)

	trigger := Trigger{Name: payload["name"]}
	triggerList.Append(trigger)
	w.WriteHeader(201)
}

func GetTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	vars := mux.Vars(r)
	trigger, err := triggerList.Get(vars["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(trigger, w)
}

func UpdateTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	vars := mux.Vars(r)
	trigger, err := triggerList.Get(vars["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	payload := unmarshal(r.Body, "cron", w)

	t := trigger.(Trigger)
	t.Schedule = payload["cron"]
	ArmTrigger(t)
	err = triggerList.Update(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteTrigger(w http.ResponseWriter, r *http.Request) {
	triggerList := GetTriggerList()

	vars := mux.Vars(r)

	triggerList.Delete(vars["trigger"])
}

func ListJobsForTrigger(w http.ResponseWriter, r *http.Request) {
	jobList := GetJobList()
	vars := mux.Vars(r)
	jobs := jobList.GetJobsWithTrigger(vars["trigger"])
	marshal(jobs, w)
}
