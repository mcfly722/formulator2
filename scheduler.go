package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
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

func (scheduler *scheduler) FinishTask(number uint64) {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	for _, task := range scheduler.Tasks {
		if task.Number == number {
			task.Done()
		}
	}

}

func (scheduler *scheduler) HTTPHandlerTask(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		task := scheduler.NewTask(r.RemoteAddr)

		json, err := json.Marshal(task)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(json)
		}
	}

	if r.Method == "DELETE" {
		vars := mux.Vars(r)
		number, err := strconv.ParseUint(vars["id"], 10, 64)
		if err == nil {
			scheduler.FinishTask(number)
		}
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
