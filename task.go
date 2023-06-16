package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Task struct {
	Number             uint64
	Sequence           string
	agent              string
	startedAt          time.Time
	lastConfirmationAt time.Time
	timeoutAt          time.Time
	timeoutSec         uint64
	done               bool
	job                func()
}

func NewTask(number uint64, sequence string, agent string, timeoutSec uint64) *Task {
	empty := func() {}

	task := Task{
		Number:             number,
		Sequence:           sequence,
		agent:              agent,
		startedAt:          time.Time{},
		lastConfirmationAt: time.Time{},
		timeoutAt:          time.Time{},
		timeoutSec:         timeoutSec,
		done:               false,
		job:                empty,
	}

	task.Reset(agent)

	return &task
}

func (task *Task) IsOutdated() bool {
	return time.Now().After(task.timeoutAt)
}

func (task *Task) IsDone() bool {
	return task.done
}

func (task *Task) Done() {
	task.done = true
}

func (task *Task) Reset(agent string) {
	now := time.Now()

	task.agent = agent

	task.startedAt = now
	task.lastConfirmationAt = now
	task.timeoutAt = now.Add(time.Duration(task.timeoutSec) * time.Second)

	fmt.Println(task)
}

func (task *Task) MarshalJSON() ([]byte, error) {
	type jsonTask struct {
		Number              uint64
		Sequence            string
		Agent               string
		StartedAt           string
		Elapsed             string
		LastConfirmationAgo string
		TimeoutedOnSec      int64
		Done                bool
	}

	j := jsonTask{
		Number:              task.Number,
		Sequence:            task.Sequence,
		Agent:               task.agent,
		StartedAt:           task.startedAt.Format("2006-01-02 15:04:05"),
		Elapsed:             time.Time{}.Add(time.Since(task.startedAt)).Format("15:04:05"),
		LastConfirmationAgo: time.Time{}.Add(time.Since(task.lastConfirmationAt)).Format("15:04:05"),
		TimeoutedOnSec:      int64(time.Since(task.lastConfirmationAt).Seconds()) - int64(task.timeoutSec),
		Done:                task.done,
	}

	return json.Marshal(j)
}

func (task *Task) SetJob(job func()) {
	task.job = job
}

func (task *Task) Do() {
	task.job()
}

func NewTaskFromServer(server string, timeoutSec uint) (*Task, error) {

	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%v%v", server, "/api/task/new"), nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var task Task

	err = json.Unmarshal(body, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (task *Task) ReportToServerWhatDone(server string, timeoutSec uint) error {
	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/api/task/%v", server, task.Number), nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	response.Body.Close()
	return nil
}
