package main

import (
	"encoding/json"
	"time"
)

type Solution struct {
	Number    uint64
	FoundedAt string
	Sequence  string
	Text      string
	Deviation float64
}

func newSolution(task *Task, deviation float64, text string) *Solution {
	return &Solution{
		Number:    task.Number,
		FoundedAt: time.Now().Format("2006-01-02 15:04:05"),
		Sequence:  task.Sequence,
		Deviation: deviation,
		Text:      text,
	}
}

func newSolutionFromString(body []byte) (*Solution, error) {
	var solution Solution

	err := json.Unmarshal(body, &solution)
	if err != nil {
		return nil, err
	}

	return &solution, nil
}
