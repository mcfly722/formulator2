package main

import (
	"encoding/json"
	"time"
)

type Solution struct {
	FoundedAt time.Time
	Text      string
	Deviation float64
}

func newSolution(task *Task, deviation float64, text string) *Solution {
	now := time.Now()

	return &Solution{
		FoundedAt: now,
		Deviation: deviation,
		Text:      text,
	}
}

func (solution *Solution) toJSONString() string {
	json, err := json.Marshal(solution)
	if err != nil {
		return ""
	}
	return string(json)
}
