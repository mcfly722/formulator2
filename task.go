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
	Number               uint64
	Sequence             string
	Solution             *Solution
	DeviationThreshold   float64
	StartedAt            time.Time
	Agent                string
	ConfirmedAt          time.Time
	ConfirmationsCounter uint64
	TimeoutAt            time.Time
	TimeoutSec           uint64
}

func newTask(number uint64, sequence string, deviationThreshold float64, agent string, timeoutSec uint64) *Task {

	now := time.Now()

	task := Task{
		Number:               number,
		Sequence:             sequence,
		Solution:             nil,
		DeviationThreshold:   deviationThreshold,
		Agent:                agent,
		StartedAt:            now,
		ConfirmedAt:          now,
		ConfirmationsCounter: 0,
		TimeoutAt:            now.Add(time.Duration(timeoutSec) * time.Second),
		TimeoutSec:           timeoutSec,
	}

	return &task
}

func (task *Task) isOutdated() bool {
	return time.Now().After(task.TimeoutAt)
}

func (task *Task) isDone() bool {
	return task.Solution != nil
}

func (task *Task) reset(agentName string) {
	now := time.Now()
	task.Agent = agentName
	task.ConfirmationsCounter = 0
	task.ConfirmedAt = now
	task.StartedAt = now
	task.TimeoutAt = now.Add(time.Duration(task.TimeoutSec) * time.Second)
}

func (task *Task) confirm() {
	now := time.Now()
	task.ConfirmedAt = now
	task.ConfirmationsCounter++
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

	fmt.Printf("-> %v\n", string(body))

	err = json.Unmarshal(body, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (task *Task) reportToServer(server string, timeoutSec uint) error {

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

func (task *Task) startWitness(solutions chan *Solution, submitIntervalSec uint, serverAddr string, agentRequestTimeoutSec uint) {

	go func(task *Task, solutions chan *Solution, submitIntervalSec uint, serverAddr string, agentRequestTimeoutSec uint) {

		ticker := time.NewTicker(time.Duration(submitIntervalSec) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				{
					// report to server task execution confirmation
					task.confirm()
					err := task.reportToServer(serverAddr, agentRequestTimeoutSec)
					if err != nil {
						fmt.Printf("whitness confirmation:%v\n", err)
					}
				}
			case solution := <-solutions:
				{
					// report to server that done
					task.Solution = solution
					err := task.reportToServer(serverAddr, agentRequestTimeoutSec)
					if err != nil {
						fmt.Printf("whitness done:%v\n", err)
					}
					return
				}
			}
		}
	}(task, solutions, submitIntervalSec, serverAddr, agentRequestTimeoutSec)

}

func (task *Task) do() *Solution {

	// sleep random pause
	time.Sleep(time.Duration(rand.Intn(12)) * time.Second)

	deviation := rand.Float64() * 10000
	//fmt.Printf("%v %v done...\n", task.Number, deviation)

	return newSolution(task, deviation, fmt.Sprintf("[%v] sequence=%v", task.Number, task.Sequence))
}
