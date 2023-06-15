package main

import (
	"encoding/json"
	"fmt"
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
}

func NewTask(number uint64, sequence string, agent string, timeoutSec uint64) *Task {

	task := Task{
		Number:     number,
		Sequence:   sequence,
		timeoutSec: timeoutSec,
	}

	task.Reset(agent)

	return &task
}

func (task *Task) IsOutdated() bool {
	return time.Now().After(task.timeoutAt)
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
	}

	j := jsonTask{
		Number:              task.Number,
		Sequence:            task.Sequence,
		Agent:               task.agent,
		StartedAt:           task.startedAt.Format("2006-01-02 15:04:05"),
		Elapsed:             time.Time{}.Add(time.Since(task.startedAt)).Format("15:04:05"),
		LastConfirmationAgo: time.Time{}.Add(time.Since(task.lastConfirmationAt)).Format("15:04:05"),
		TimeoutedOnSec:      int64(time.Since(task.lastConfirmationAt).Seconds()) - int64(task.timeoutSec),
	}

	return json.Marshal(j)
}
