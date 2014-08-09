package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	. "github.com/jakecoffman/gorunner/service"
	"github.com/nu7hatch/gouuid"
)

var nothing = map[string]string{}

func errHelp(msg string) map[string]interface{} {
	return map[string]interface{}{"error": msg}
}

// General

func app(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/app.html")
}

func wsHandler(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	// Upgrade the HTTP connection to a websocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		return http.StatusBadRequest, errHelp("Not a websocket handshake")
	} else if err != nil {
		return http.StatusInternalServerError, errHelp(err.Error())
	}
	conn := NewConnection(ws)
	c.Hub().Register(conn)
	defer c.Hub().Unregister(conn)
	go conn.Writer()
	conn.Reader()
	return http.StatusOK, nothing
}

// Jobs

func listJobs(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.JobList().Dump()
}

func addJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)

	err := c.JobList().Append(Job{Name: payload["name"], Status: "New"})
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusCreated, nothing
}

func getJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	return http.StatusOK, job
}

func deleteJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	err = c.JobList().Delete(job.ID())
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, nothing
}

func addTaskToJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	payload := unmarshal(r.Body, "task", w)
	j.AppendTask(payload["task"])
	c.JobList().Update(j)

	return http.StatusCreated, nothing
}

func removeTaskFromJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	taskPosition, err := strconv.Atoi(vars["task"])
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	j.DeleteTask(taskPosition)
	c.JobList().Update(j)
	return http.StatusOK, nothing
}

func addTriggerToJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	payload := unmarshal(r.Body, "trigger", w)

	j.AppendTrigger(payload["trigger"])
	t, err := c.TriggerList().Get(payload["trigger"])
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	c.Executor().ArmTrigger(t.(Trigger))
	c.JobList().Update(j)

	return http.StatusCreated, nothing
}

func removeTriggerFromJob(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	job, err := c.JobList().Get(vars["job"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	j := job.(Job)

	t := vars["trigger"]
	j.DeleteTrigger(t)
	c.JobList().Update(j)

	// If Trigger is no longer attached to any Jobs, remove it from Cron to save cycles
	jobs := c.JobList().GetJobsWithTrigger(t)

	if len(jobs) == 0 {
		c.Executor().DisarmTrigger(t)
	}
	return http.StatusOK, nothing
}

// Run

func listRuns(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
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

	return http.StatusOK, c.RunList().GetRecent(o, l)
}

func addRun(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "job", w)

	job, err := c.JobList().Get(payload["job"])
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
		task, err := c.TaskList().Get(taskName)
		if err != nil {
			panic(err)
		}
		t := task.(Task)
		tasks = append(tasks, t)
	}
	err = c.RunList().AddRun(id.String(), j, tasks)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusCreated, map[string]string{"uuid": id.String()}
}

func getRun(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	run, err := c.RunList().Get(vars["run"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusOK, run
}

// Tasks

func listTasks(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.TaskList().Dump()
}

func addTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)

	err := c.TaskList().Append(Task{payload["name"], ""})
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	return http.StatusCreated, nothing
}

func getTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.TaskList().Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusOK, task
}

func updateTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.TaskList().Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	payload := unmarshal(r.Body, "script", w)
	t := task.(Task)
	t.Script = payload["script"]
	c.TaskList().Update(t)
	return http.StatusOK, nothing
}

func deleteTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	task, err := c.TaskList().Get(vars["task"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	c.TaskList().Delete(task.ID())
	return http.StatusOK, nothing
}

func listJobsForTask(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	jobs := c.JobList().GetJobsWithTask(vars["task"])
	return http.StatusOK, jobs
}

// Triggers

func listTriggers(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	return http.StatusOK, c.TriggerList().Dump()
}

func addTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	payload := unmarshal(r.Body, "name", w)
	trigger := Trigger{Name: payload["name"]}
	c.TriggerList().Append(trigger)
	return http.StatusCreated, nothing
}

func getTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	trigger, err := c.TriggerList().Get(vars["trigger"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}
	return http.StatusNotFound, trigger
}

func updateTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	trigger, err := c.TriggerList().Get(vars["trigger"])
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	payload := unmarshal(r.Body, "cron", w)

	t := trigger.(Trigger)
	t.Schedule = payload["cron"]
	c.Executor().ArmTrigger(t)
	err = c.TriggerList().Update(t)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, nothing
}

func deleteTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	c.TriggerList().Delete(vars["trigger"])
	return http.StatusOK, nothing
}

func listJobsForTrigger(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	vars := mux.Vars(r)
	jobs := c.JobList().GetJobsWithTrigger(vars["trigger"])
	return http.StatusOK, jobs
}
