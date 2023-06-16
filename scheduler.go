package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type scheduler struct {
	state *State
	Tasks []*Task

	lastCounter  uint64
	lastSequence string

	defaultTimeoutSec uint64
	ready             sync.Mutex
}

func NewScheduler(state *State, defaultTimeoutSec uint64) *scheduler {

	return &scheduler{
		state:             state,
		Tasks:             []*Task{},
		lastCounter:       state.Counter,
		lastSequence:      state.LastSequence,
		defaultTimeoutSec: defaultTimeoutSec,
	}
}

func (scheduler *scheduler) NewTask(agent string) Task {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	// search and return first outdated
	for _, task := range scheduler.Tasks {
		if task.IsOutdated() && !task.IsDone() {
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

func (scheduler *scheduler) FinishTask(number uint64, solution *Solution) {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	for _, task := range scheduler.Tasks {
		if task.Number == number {
			task.Done()
		}
	}

	// take and report about first done tasks
	for len(scheduler.Tasks) > 0 && (scheduler.Tasks[0].done) {
		task := scheduler.Tasks[0]

		scheduler.state.ReportAboutSolution(task, solution)

		scheduler.Tasks = scheduler.Tasks[1:] // remove first done task
	}

}

func (scheduler *scheduler) HTTPHandlerTask(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		task := scheduler.NewTask(r.RemoteAddr)

		json, err := json.Marshal(task)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(json)
	}

	if r.Method == "DELETE" {

		vars := mux.Vars(r)
		number, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		solution, err := NewSolutionFromString(body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		scheduler.FinishTask(number, solution)
	}
}

func (scheduler *scheduler) HTTPHandlerTasks(w http.ResponseWriter, r *http.Request) {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	json, err := json.Marshal(scheduler.Tasks)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(json)
	}

}
