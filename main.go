package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	// server parameters
	stateFile := flag.String("stateFile", "state.json", "file for current state")
	stateSaveIntervalSec := flag.Uint("stateSaveIntervalSec", 5, "state save interval in secs")
	saveFirstNBestSolutions := flag.Uint("saveFirstNBestSolutions", 15, "save first N solutions")
	samplesFile := flag.String("samplesFile", "samples\\exponent\\exponent.json", "file of points for required function")
	listenAddr := flag.String("listenAddr", "localhost:8080", "bind address for server")
	taskTimeoutSec := flag.Uint64("taskTimeoutSec", 5, "tasks timeout in secs")

	// agent parameters
	serverAddr := flag.String("serverAddr", "", "server address")
	agentName := flag.String("agentName", "", "agent name")
	agentThreads := flag.Uint("agentThreads", 5, "number of agent threads")
	agentErrorSleepSec := flag.Uint("agentErrorSleep", 1, "seconds before next try after error")
	agentRequestTimeoutSec := flag.Uint("agentRequestTimeoutSec", 5, "timeout for agent->server request")

	flag.Parse()

	if len(*serverAddr) == 0 {
		var state *State

		// read state or create new
		if isStateFileExist(*stateFile) {
			loadedState, err := loadStateFromFile(*stateFile, *saveFirstNBestSolutions)
			if err != nil {
				panic(err)
			}
			state = loadedState

			fmt.Printf("previous state founded and loaded from %v\n", *stateFile)

		} else {
			// new state
			points, err := loadPointsFromFile(*samplesFile)
			if err != nil {
				panic(err)
			}

			state = NewState(points, *stateFile, *saveFirstNBestSolutions)
			fmt.Printf("new state created. Loaded %v sample points from %v\n", len(state.Points), *stateFile)

		}

		state.startRegularSaving(*stateSaveIntervalSec)

		server := newServer(*listenAddr, newScheduler(state, *taskTimeoutSec), state)

		fmt.Printf("starting server on %v\n", *listenAddr)

		if err := server.start(*listenAddr); err != nil {
			log.Fatal(err)
		}

	} else {
		if len(*agentName) == 0 {
			hostname, err := os.Hostname()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			*agentName = hostname
		}

		fmt.Printf("starting agent\n")

		rand.Seed(time.Now().UnixNano())

		pool := newPool(*agentName, *agentThreads, *serverAddr, *agentErrorSleepSec, *agentRequestTimeoutSec)
		pool.Start()
	}
}
