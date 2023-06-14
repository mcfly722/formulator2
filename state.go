package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

type State struct {
	Points   []Point
	LastTree string
	Counter  uint64

	ready sync.Mutex
}

func NewState(points []Point) *State {
	return &State{
		Points:   points,
		LastTree: "",
		Counter:  0,
	}
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
	} else {
		w.Write(json)
	}
}
