package main

import (
	"encoding/json"
	"fmt"
	"math"
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
	saveFirstNBestSolutions uint
	Counter                 uint64
	LastSequence            string
	ready                   sync.Mutex
}

func NewState(points []Point, fileName string, saveFirstNBestSolutions uint) *State {
	return &State{
		FileName:                fileName,
		Points:                  points,
		Solutions:               []Solution{},
		saveFirstNBestSolutions: saveFirstNBestSolutions,
		Counter:                 0,
	}
}

func (state *State) getDeviationThreshold() float64 {
	state.ready.Lock()
	defer state.ready.Unlock()

	// if there are no solutions found yet, return max possible
	if len(state.Solutions) == 0 {
		return math.MaxFloat64
	}

	// return max deviation from all registered solutions
	i := len(state.Solutions) - 1
	return state.Solutions[i].Deviation
}

func (state *State) startRegularSaving(intervalSec uint) {
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

func (state *State) reportAboutSolution(task *Task) {
	state.ready.Lock()
	defer state.ready.Unlock()

	fmt.Printf("[%v] %v\n", task.Number, task.Solution.toJSONString())

	state.Counter++
	state.LastSequence = task.Sequence

	// not optimized insert in to sorted slice with right shift
	state.Solutions = append(state.Solutions, *task.Solution)

	sort.Slice(state.Solutions, func(i, j int) bool {
		return state.Solutions[i].Deviation < state.Solutions[j].Deviation
	})

	// trim only first N elements from solutions slice
	l := len(state.Solutions)
	if l > int(state.saveFirstNBestSolutions) {
		l = int(state.saveFirstNBestSolutions)
	}
	state.Solutions = state.Solutions[:l]
}

func isStateFileExist(stateFile string) bool {
	if _, err := os.Stat(stateFile); err == nil {
		return true
	}

	return false
}

func loadStateFromFile(stateFile string, saveFirstNBestSolutions uint) (*State, error) {
	body, err := os.ReadFile(stateFile)

	if err != nil {
		return nil, err
	}

	state := State{}

	err = json.Unmarshal(body, &state)
	if err != nil {
		return nil, err
	}

	state.saveFirstNBestSolutions = saveFirstNBestSolutions
	return &state, nil
}

func (state *State) httpHandlerPoints(w http.ResponseWriter, r *http.Request) {
	state.ready.Lock()
	defer state.ready.Unlock()

	json, err := json.Marshal(state.Points)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(json)
}

func (state *State) httpHandlerState(w http.ResponseWriter, r *http.Request) {
	state.ready.Lock()
	defer state.ready.Unlock()

	json, err := json.Marshal(state)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(json)
}
