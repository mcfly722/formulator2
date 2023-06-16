package main

import (
	"bytes"
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
	finishedAt         time.Time
	finished           bool
	job                func()
}

func newTask(number uint64, sequence string, agent string, timeoutSec uint64) *Task {
	empty := func() {}

	task := Task{
		Number:             number,
		Sequence:           sequence,
		agent:              agent,
		startedAt:          time.Time{},
		lastConfirmationAt: time.Time{},
		timeoutAt:          time.Time{},
		timeoutSec:         timeoutSec,
		finishedAt:         time.Time{},
		finished:           false,
		job:                empty,
	}

	task.reset(agent)

	return &task
}

func (task *Task) isOutdated() bool {
	return time.Now().After(task.timeoutAt)
}

func (task *Task) isDone() bool {
	return task.finished
}

func (task *Task) done() {
	task.finishedAt = time.Now()
	task.finished = true
}

func (task *Task) reset(agent string) {
	now := time.Now()

	task.agent = agent

	task.startedAt = now
	task.lastConfirmationAt = now
	task.timeoutAt = now.Add(time.Duration(task.timeoutSec) * time.Second)
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

	elapsed := time.Time{}.Add(time.Since(task.startedAt)).Format("15:04:05")
	if task.isDone() {
		elapsed = time.Time{}.Add(task.finishedAt.Sub(task.startedAt)).Format("15:04:05")
	}

	j := jsonTask{
		Number:              task.Number,
		Sequence:            task.Sequence,
		Agent:               task.agent,
		StartedAt:           task.startedAt.Format("2006-01-02 15:04:05"),
		Elapsed:             elapsed,
		LastConfirmationAgo: time.Time{}.Add(time.Since(task.lastConfirmationAt)).Format("15:04:05"),
		TimeoutedOnSec:      int64(time.Since(task.lastConfirmationAt).Seconds()) - int64(task.timeoutSec),
		Done:                task.isDone(),
	}

	return json.Marshal(j)
}

func (task *Task) setJob(job func()) {
	task.job = job
}

func (task *Task) do() {
	task.job()
}

func newTaskFromServer(server string, timeoutSec uint) (*Task, error) {

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

func (task *Task) reportToServerWhatDone(server string, timeoutSec uint, solution *Solution) error {

	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	solutionJSON, err := json.Marshal(solution)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/api/task/%v", server, task.Number), bytes.NewBuffer(solutionJSON))
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
