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
	state             *State
	tasks             []*Task
	lastCounter       uint64
	lastSequence      string
	defaultTimeoutSec uint64
	ready             sync.Mutex
}

func newScheduler(state *State, defaultTimeoutSec uint64) *scheduler {

	return &scheduler{
		state:             state,
		tasks:             []*Task{},
		lastCounter:       state.Counter,
		lastSequence:      state.LastSequence,
		defaultTimeoutSec: defaultTimeoutSec,
	}
}

func (scheduler *scheduler) newTask(agent string) Task {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	// search and return first outdated
	for _, task := range scheduler.tasks {
		if task.isOutdated() && !task.isDone() {
			task.reset(agent)
			return *task
		}
	}

	// if there are no outdated tasks, create new one

	scheduler.lastCounter++

	lastTree := fmt.Sprintf("[%v]", scheduler.lastCounter)

	newTask := newTask(scheduler.lastCounter, lastTree, agent, scheduler.defaultTimeoutSec)
	scheduler.tasks = append(scheduler.tasks, newTask)

	return *newTask
}

func (scheduler *scheduler) finishTask(number uint64, solution *Solution) {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	for _, task := range scheduler.tasks {
		if task.Number == number {
			task.done()
		}
	}

	// take and report about first done tasks
	for len(scheduler.tasks) > 0 && (scheduler.tasks[0].isDone()) {
		task := scheduler.tasks[0]

		scheduler.state.reportAboutSolution(task, solution)

		scheduler.tasks = scheduler.tasks[1:] // remove first done task
	}

}

func (scheduler *scheduler) httpHandlerTask(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		task := scheduler.newTask(r.RemoteAddr)

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

		solution, err := newSolutionFromString(body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		scheduler.finishTask(number, solution)
	}
}

func (scheduler *scheduler) httpHandlerTasks(w http.ResponseWriter, r *http.Request) {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	json, err := json.Marshal(scheduler.tasks)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(json)
	}

}
