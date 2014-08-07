package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
)

var nothing = map[string]string{}

func errHelp(msg string) map[string]interface{} {
	return map[string]interface{}{"error": msg}
}

// General

func App(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/app.html")
}

func WsHandler(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	// Upgrade the HTTP connection to a websocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		return http.StatusBadRequest, errHelp("Not a websocket handshake")
	} else if err != nil {
		return http.StatusInternalServerError, errHelp(err.Error())
	}
	context := NewConnection(ws)
	c.hub.Register(context)
	defer c.hub.Unregister(context)
	go context.Writer()
	context.Reader()
	return http.StatusOK, nothing
}

// Jobs

func ListJobs(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.jobList.Dump()
}

func AddJob(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)

	err := c.jobList.Append(Job{Name: payload["name"], Status: "New"})
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusCreated, nothing
}

func GetJob(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.jobList.Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	return http.StatusOK, job
}

func DeleteJob(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.jobList.Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	err = c.jobList.Delete(job.ID())
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, nothing
}

func AddTaskToJob(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.jobList.Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	payload := unmarshal(r.Body, "task", w)
	j.AppendTask(payload["task"])
	c.jobList.Update(j)

	return http.StatusCreated, nothing
}

func RemoveTaskFromJob(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.jobList.Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	taskPosition, err := strconv.Atoi(vars["task"])
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	j.DeleteTask(taskPosition)
	c.jobList.Update(j)
	return http.StatusOK, nothing
}

func AddTriggerToJob(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.jobList.Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	payload := unmarshal(r.Body, "trigger", w)

	j.AppendTrigger(payload["trigger"])
	t, err := c.triggerList.Get(payload["trigger"])
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	c.executor.ArmTrigger(t.(Trigger))
	c.jobList.Update(j)

	return http.StatusCreated, nothing
}

func RemoveTriggerFromJob(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.jobList.Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	t := vars["trigger"]
	j.DeleteTrigger(t)
	c.jobList.Update(j)

	// If Trigger is no longer attached to any Jobs, remove it from Cron to save cycles
	jobs := c.jobList.GetJobsWithTrigger(t)

	if len(jobs) == 0 {
		c.executor.DisarmTrigger(t)
	}
	return http.StatusOK, nothing
}

// Run

func ListRuns(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
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
		return http.StatusBadRequest, err.Error()
	}

	l, err := strconv.Atoi(length)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	return http.StatusOK, c.runList.GetRecent(o, l)
}

func AddRun(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "job", w)

	job, err := c.jobList.Get(payload["job"])
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	j := job.(Job)

	id, err := uuid.NewV4()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	var tasks []Task
	for _, taskName := range j.Tasks {
		task, err := c.taskList.Get(taskName)
		if err != nil {
			panic(err)
		}
		t := task.(Task)
		tasks = append(tasks, t)
	}
	err = c.runList.AddRun(id.String(), j, tasks)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusCreated, map[string]string{"uuid": id.String()}
}

func GetRun(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	run, err := c.runList.Get(vars["run"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusOK, run
}

// Tasks

func ListTasks(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.taskList.Dump()
}

func AddTask(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)

	err := c.taskList.Append(Task{payload["name"], ""})
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	return http.StatusCreated, nothing
}

func GetTask(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.taskList.Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusOK, task
}

func UpdateTask(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.taskList.Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	payload := unmarshal(r.Body, "script", w)
	t := task.(Task)
	t.Script = payload["script"]
	c.taskList.Update(t)
	return http.StatusOK, nothing
}

func DeleteTask(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.taskList.Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	c.taskList.Delete(task.ID())
	return http.StatusOK, nothing
}

func ListJobsForTask(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	jobs := c.jobList.GetJobsWithTask(vars["task"])
	return http.StatusOK, jobs
}

// Triggers

func ListTriggers(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.triggerList.Dump()
}

func AddTrigger(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)
	trigger := Trigger{Name: payload["name"]}
	c.triggerList.Append(trigger)
	return http.StatusCreated, nothing
}

func GetTrigger(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	trigger, err := c.triggerList.Get(vars["trigger"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusNotFound, trigger
}

func UpdateTrigger(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	trigger, err := c.triggerList.Get(vars["trigger"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	payload := unmarshal(r.Body, "cron", w)

	t := trigger.(Trigger)
	t.Schedule = payload["cron"]
	c.executor.ArmTrigger(t)
	err = c.triggerList.Update(t)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, nothing
}

func DeleteTrigger(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	c.triggerList.Delete(vars["trigger"])
	return http.StatusOK, nothing
}

func ListJobsForTrigger(c *context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	jobs := c.jobList.GetJobsWithTrigger(vars["trigger"])
	return http.StatusOK, jobs
}
