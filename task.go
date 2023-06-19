package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Task struct {
	Number             uint64
	Sequence           string
	Solution           *Solution
	StartedAt          time.Time
	Agent              string
	lastConfirmationAt time.Time
	TimeoutAt          time.Time
	TimeoutSec         uint64
}

func newTask(number uint64, sequence string, agent string, timeoutSec uint64) *Task {

	now := time.Now()

	task := Task{
		Number:             number,
		Sequence:           sequence,
		Solution:           nil,
		Agent:              agent,
		StartedAt:          now,
		lastConfirmationAt: now,
		TimeoutAt:          now.Add(time.Duration(timeoutSec) * time.Second),
		TimeoutSec:         timeoutSec,
	}

	task.reset(agent)

	return &task
}

func (task *Task) isOutdated() bool {
	return time.Now().After(task.TimeoutAt)
}

func (task *Task) isDone() bool {
	return task.Solution != nil
}

func (task *Task) reset(agent string) {
	now := time.Now()

	task.Agent = agent

	task.StartedAt = now
	task.lastConfirmationAt = now
	task.TimeoutAt = now.Add(time.Duration(task.TimeoutSec) * time.Second)
}

func newTaskFromServer(server string, agentName string, timeoutSec uint) (*Task, error) {
	var task Task

	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%v/api/task/%v", server, agentName), nil)
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

	//fmt.Printf("-> %v\n", string(body))

	err = json.Unmarshal(body, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (task *Task) reportToServerWhatDone(server string, timeoutSec uint) error {

	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%v/api/task/%v", server, task.Agent), bytes.NewBuffer(taskJSON))
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	response.Body.Close()

	fmt.Printf("<- %v\n", string(taskJSON))
	return nil
}

func (task *Task) do() {
	// sleep random pause
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	deviation := rand.Float64() * 1000
	//fmt.Printf("%v %v done...\n", task.Number, deviation)

	task.Solution = newSolution(task, deviation, fmt.Sprintf("[%v] sequence=%v", task.Number, task.Sequence))
}
