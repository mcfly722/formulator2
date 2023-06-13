package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	stateFile := *flag.String("stateFile", "state.json", "file for current state")
	samplesFile := *flag.String("samplesFile", "samples\\exponent\\exponent.json", "file of points for required function")

	listenAddr := *flag.String("listenAddr", "localhost:8080", "bind address for server")
	serverAddr := *flag.String("serverAddr", "", "server address")

	defaultTimeoutSec := *flag.Uint64("timeoutSec", 5, "tasks timeout in secs")

	flag.Parse()

	if len(serverAddr) == 0 {
		// server

		var state *State

		// read state or create new
		if IsStateFileExist(stateFile) {
			loadedState, err := LoadStateFromFile(stateFile)
			if err != nil {
				panic(err)
			}
			state = loadedState

			fmt.Printf("previous state founded and loaded from %v", stateFile)

		} else {
			// new state
			points, err := LoadPointsFromFile(samplesFile)
			if err != nil {
				panic(err)
			}

			state = NewState(points)
			fmt.Printf("new state created. Loaded %v sample points from %v\n", len(state.Points), stateFile)
		}

		server := NewServer(listenAddr, NewScheduler(state, defaultTimeoutSec))

		fmt.Printf("starting server on %v\n", listenAddr)

		if err := server.start(listenAddr); err != nil {
			log.Fatal(err)
		}

	}

}
