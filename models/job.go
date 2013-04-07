package models

type Job struct {
	Name  string
	Tasks []string
}

func (j *Job) Append(task string) {
	j.Tasks = append(j.Tasks, task)
}

func (j *Job) Delete(taskPosition int) error {
	i := taskPosition
	j.Tasks = j.Tasks[:i+copy(j.Tasks[i:], j.Tasks[i+1:])]
	return nil
}
