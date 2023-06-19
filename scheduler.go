package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	lastSequence := fmt.Sprintf("[%v]", scheduler.lastCounter)
	newTask := newTask(scheduler.lastCounter, lastSequence, agent, scheduler.defaultTimeoutSec)
	scheduler.tasks = append(scheduler.tasks, newTask)

	return *newTask
}

func (scheduler *scheduler) finishTask(task *Task) {

	scheduler.ready.Lock()
	defer scheduler.ready.Unlock()

	// mark task as done
	for j := 0; j < len(scheduler.tasks); j++ {
		if scheduler.tasks[j].Number == task.Number {
			scheduler.tasks[j] = task
		}
	}

	/*
		fmt.Printf("BEFORE:")
		for _, task := range scheduler.tasks {
			fmt.Printf("%v[%v],", task.Number, task.isDone())
		}
		fmt.Printf("\n")
	*/

	// take and report about first done tasks
	var i int = 0
	for i = 0; i < len(scheduler.tasks) && scheduler.tasks[i].isDone(); i++ {
		scheduler.state.reportAboutSolution(scheduler.tasks[i])
	}

	// remove first i done tasks
	scheduler.tasks = scheduler.tasks[i:]

	/*
		fmt.Printf(" AFTER:")
		for _, task := range scheduler.tasks {
			fmt.Printf("%v[%v],", task.Number, task.isDone())
		}
		fmt.Printf("\n")
	*/

}

func (scheduler *scheduler) httpHandlerTask(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		//agentName
		vars := mux.Vars(r)
		if agentName, found := vars["agentName"]; !found {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("agent name is not specified"))
			return
		} else {
			task := scheduler.newTask(agentName)

			json, err := json.Marshal(task)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(json)
		}

	}

	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		var task Task

		err = json.Unmarshal(body, &task)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		//fmt.Printf("->%v", string(body))

		scheduler.finishTask(&task)
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
