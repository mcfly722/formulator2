package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type scheduler struct {
	state *State
	Tasks []*Task

	lastCounter uint64
	lastTree    string

	defaultTimeoutSec uint64
	ready             sync.Mutex
}

func NewScheduler(state *State, defaultTimeoutSec uint64) *scheduler {
	return &scheduler{
		state:             state,
		Tasks:             []*Task{},
		lastCounter:       state.Counter,
		lastTree:          state.LastTree,
		defaultTimeoutSec: defaultTimeoutSec,
	}
}

func (scheduler *scheduler) NewTask(agent string) Task {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	// search and return first outdated
	for _, task := range scheduler.Tasks {
		if task.IsOutdated() {
			task.Reset(agent)
			return *task
		}
	}

	// if there are no outdated tasks, create new one

	scheduler.lastCounter++

	lastTree := fmt.Sprintf("[%v]", scheduler.lastCounter)

	newTask := NewTask(scheduler.lastCounter, lastTree, agent, scheduler.defaultTimeoutSec)
	scheduler.Tasks = append(scheduler.Tasks, newTask)

	return *newTask
}

func (scheduler *scheduler) HTTPHandlerNewTask(w http.ResponseWriter, r *http.Request) {
	task := scheduler.NewTask(r.RemoteAddr)

	json, err := json.Marshal(task)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(json)
	}
}

func (scheduler *scheduler) HTTPHandlerListTasks(w http.ResponseWriter, r *http.Request) {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	json, err := json.Marshal(scheduler.Tasks)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(json)
	}
}
