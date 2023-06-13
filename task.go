package main

import (
	"fmt"
	"time"
)

type Task struct {
	Number             uint64
	Sequence           string
	Agent              string
	StartedAt          time.Time
	LastConfirmationAt time.Time
	TimeoutAt          time.Time
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
	return time.Now().After(task.TimeoutAt)
}

func (task *Task) Reset(agent string) {
	now := time.Now()

	task.Agent = agent

	task.StartedAt = now
	task.LastConfirmationAt = now
	task.TimeoutAt = now.Add(time.Duration(task.timeoutSec) * time.Second)

	fmt.Println(task)
}
