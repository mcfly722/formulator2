package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type State struct {
	FileName                string
	Points                  []Point
	Solutions               []Solution
	SaveFirstNBestSolutions uint
	Counter                 uint64
	LastSequence            string
	ready                   sync.Mutex
}

func NewState(points []Point, fileName string, saveFirstNBestSolutions uint) *State {
	return &State{
		FileName:                fileName,
		Points:                  points,
		Solutions:               []Solution{},
		SaveFirstNBestSolutions: saveFirstNBestSolutions,
		Counter:                 0,
	}
}

func (state *State) StartRegularSaving(intervalSec uint) {
	go func() {

		for {
			time.Sleep(time.Duration(intervalSec) * time.Second)

			fmt.Printf("saving state...")

			state.ready.Lock()
			json, err := json.Marshal(state)
			state.ready.Unlock()

			if err != nil {
				fmt.Printf("%v\n", err.Error())
			} else {

				f, _ := os.OpenFile(state.FileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

				f.Write(json)

				fmt.Printf("ok\n")
				f.Close()
			}
		}
	}()
}

func (state *State) ReportAboutSolution(task *Task, solution *Solution) {
	state.ready.Lock()
	defer state.ready.Unlock()

	state.Counter++
	state.LastSequence = task.Sequence

	// not optimized insert in to sorted slice with right shift
	state.Solutions = append(state.Solutions, *solution)

	sort.Slice(state.Solutions, func(i, j int) bool {
		return state.Solutions[i].Deviation < state.Solutions[j].Deviation
	})

	// trim only first N elements from solutions slice
	l := len(state.Solutions)
	if l > int(state.SaveFirstNBestSolutions) {
		l = int(state.SaveFirstNBestSolutions)
	}
	state.Solutions = state.Solutions[:l]
}

func IsStateFileExist(stateFile string) bool {
	if _, err := os.Stat(stateFile); err == nil {
		return true
	}

	return false
}

func LoadStateFromFile(stateFile string) (*State, error) {
	body, err := os.ReadFile(stateFile)

	if err != nil {
		return nil, err
	}

	state := State{}

	err = json.Unmarshal(body, &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (state *State) HTTPHandlerPoints(w http.ResponseWriter, r *http.Request) {
	state.ready.Lock()
	defer state.ready.Unlock()

	json, err := json.Marshal(state.Points)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(json)
}

func (state *State) HTTPHandlerState(w http.ResponseWriter, r *http.Request) {
	state.ready.Lock()
	defer state.ready.Unlock()

	json, err := json.Marshal(state)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(json)
}

func (state *State) HTTPHandlerSolutions(w http.ResponseWriter, r *http.Request) {
	state.ready.Lock()
	defer state.ready.Unlock()

	json, err := json.Marshal(state.Solutions)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(json)
}
