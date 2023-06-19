package main

import (
	"time"
)

type Solution struct {
	FoundedAt time.Time
	Elapsed   string
	Text      string
	Deviation float64
}

func newSolution(task *Task, deviation float64, text string) *Solution {
	now := time.Now()

	return &Solution{
		FoundedAt: now,
		Elapsed:   time.Time{}.Add(now.Sub(task.StartedAt)).Format("15:04:05"),
		Deviation: deviation,
		Text:      text,
	}
}
