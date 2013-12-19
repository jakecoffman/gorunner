package models

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"time"
	"bufio"
	"io"
	"log"
)

type Result struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Task   Task      `json:"task"`
	Output string    `json:"output"`
	Error  string    `json:"error"`
}

type Run struct {
	UUID    string    `json:"uuid"`
	Job     Job       `json:"job"`
	Tasks   []Task    `json:"tasks"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Results []*Result  `json:"results"`
	Status  string    `json:"status"`
}

func (r Run) ID() string {
	return r.UUID
}

type RunList struct {
	list
}

func (l *RunList) Load() {
	bytes := readFile(l.fileName)
	var runs []Run
	err := json.Unmarshal([]byte(string(bytes)), &runs)
	if err != nil {
		panic(err)
	}
	l.elements = nil
	for _, run := range runs {
		l.elements = append(l.elements, run)
	}
}

func (j *RunList) Len() int {
	j.RLock()
	defer j.RUnlock()

	return len(j.elements)
}

func (l *RunList) Less(i, j int) bool {
	l.RLock()
	defer l.RUnlock()

	return l.elements[i].(Run).Start.Before(l.elements[j].(Run).Start)
}

func (l *RunList) Swap(i, j int) {
	l.RLock()
	defer l.RUnlock()

	l.elements[i], l.elements[j] = l.elements[j], l.elements[i]
}

func (l RunList) GetRecent(offset, length int) []elementer {
	runs := l.elements
	if offset != -1 {
		if offset >= len(runs) {
			return nil
		}
		if length != -1 && offset+length < len(runs) {
			runs = runs[offset : offset+length]
		} else {
			runs = runs[offset:]
		}
	} else {
		if length != -1 {
			runs = runs[:length]
		}
	}
	return runs
}

func (j *RunList) AddRun(UUID string, job Job, tasks []Task) error {
	run := Run{UUID: UUID, Job: job, Tasks: tasks, Start: time.Now(), Status: "New"}
	var found bool = false
	for _, j := range j.elements {
		if run.UUID == j.(Run).UUID {
			found = true
		}
	}
	if found {
		return errors.New("Run with that name found in list")
	}
	j.Lock()
	defer j.Unlock()

	j.elements = append(j.elements, run)
	go j.execute(&run)
	j.save()
	return nil
}

func (l *RunList) execute(r *Run) {
	r.Status = "Running"
	for _, task := range r.Tasks {
		result := &Result{Start: time.Now(), Task: task}
		r.Results = append(r.Results, result)
		l.Update(*r)
		shell, commandArg := getShell()
		cmd := exec.Command(shell, commandArg, task.Script)
		//out, err := cmd.Output()
		outPipe, err := cmd.StdoutPipe()
		if err != nil {
			reportRunError(l, r, result, err)
			return
		}
		errPipe, err := cmd.StderrPipe()
		if err != nil {
			outPipe.Close()
			reportRunError(l, r, result, err)
			return
		}

		go consumeOutput(outPipe, result)
		go consumeOutput(errPipe, result)

		if err := cmd.Start(); err != nil {
			reportRunError(l, r, result, err)
			return
		}

		if err := cmd.Wait(); err != nil {
			reportRunError(l, r, result, err)
			return
		}
		result.End = time.Now()
		if err != nil {
			reportRunError(l, r, result, err)
			return
		}
		l.Update(*r)
	}
	r.End = time.Now()
	r.Status = "Done"
	jobList := GetJobList()
	job, err := jobList.Get(r.Job.Name)
	if err != nil {
		return
	}
	j := job.(Job)
	j.Status = "Ok" 
	jobList.Update(job)
	l.Update(*r)
}

func consumeOutput(reader io.ReadCloser, result *Result) {
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		//log.Println("   Consumed: ",  text)
		result.Output += text + "\n"
	}
	if err := scanner.Err(); err != nil {
		log.Println("reading standard input:", err)
	}
}

func reportRunError(l *RunList, r *Run, result *Result, err error) {
	log.Println("Reporting error", err)
	result.Error = err.Error()
	r.Status = "Failed"
	r.End = time.Now()
	l.Update(*r)
	jobList := GetJobList()
	job, err := jobList.Get(r.Job.Name)
	if err != nil {
		return
	}
	j := job.(Job)
	j.Status = "Failing"
	jobList.Update(job)
	return
}

func getShell() (string, string) {
	var shell = os.Getenv("SHELL")
	if ("" != shell ) {
		return shell,"-c"
	}
	return "cmd", "/C"
}
