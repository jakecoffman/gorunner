package models

import (
	"io/ioutil"
	"os"
)

const (
	jobsFile = "jobs.json"
	runsFile = "runs.json"
	tasksFile = "tasks.json"
	triggersFile = "triggers.json"
)

var (
	jobList JobList
	taskList TaskList
	runList RunList
	triggerList TriggerList
)

func init() {
	Load(&jobList, jobsFile)
	Load(&runList, runsFile)
	Load(&taskList, tasksFile)
	Load(&triggerList, triggersFile)
}

func GetJobList() *JobList {
	return &jobList
}

func GetRunList() *RunList {
	return &runList
}

func GetTaskList() *TaskList {
	return &taskList
}

func GetTriggerList() *TriggerList {
	return &triggerList
}

type Serializable interface{
	dumps() string
	loads(s string)
}

func Save(list Serializable, filePath string) {
	bytes := list.dumps()
	writeFile([]byte(bytes), filePath)
}

func Load(list Serializable, filePath string) {
	bytes := readFile(filePath)
	list.loads(string(bytes))
}

func writeFile(bytes []byte, filePath string) {
	err := ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func readFile(filePath string) []byte {
	_, err := os.Stat(filePath)
	if err != nil {
		println("Couldn't file file, creating fresh:", filePath)
		err = ioutil.WriteFile(filePath, []byte("[]"), 0644)
		if err != nil {
			panic(err)
		}
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return bytes
}
