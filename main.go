package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {

	// server parameters
	stateFile := flag.String("stateFile", "state.json", "file for current state")
	stateSaveIntervalSec := flag.Uint("stateSaveIntervalSec", 5, "state save interval in secs")
	saveFirstNBestSolutions := flag.Uint("saveFirstNBestSolutions", 5, "save first N solutions")
	samplesFile := flag.String("samplesFile", "samples\\exponent\\exponent.json", "file of points for required function")
	listenAddr := flag.String("listenAddr", "localhost:8080", "bind address for server")
	taskTimeoutSec := flag.Uint64("taskTimeoutSec", 5, "tasks timeout in secs")

	// agent parameters
	serverAddr := flag.String("serverAddr", "", "server address")
	agentThreads := flag.Uint("agentThreads", 5, "number of agent threads")
	agentErrorSleepSec := flag.Uint("agentErrorSleep", 1, "seconds before next try after error")
	agentRequestTimeoutSec := flag.Uint("agentRequestTimeoutSec", 5, "timeout for agent->server request")

	flag.Parse()

	if len(*serverAddr) == 0 {
		var state *State

		// read state or create new
		if IsStateFileExist(*stateFile) {
			loadedState, err := LoadStateFromFile(*stateFile)
			if err != nil {
				panic(err)
			}
			state = loadedState

			fmt.Printf("previous state founded and loaded from %v\n", *stateFile)

		} else {
			// new state
			points, err := LoadPointsFromFile(*samplesFile)
			if err != nil {
				panic(err)
			}

			state = NewState(points, *stateFile, *saveFirstNBestSolutions)
			fmt.Printf("new state created. Loaded %v sample points from %v\n", len(state.Points), *stateFile)

		}

		state.StartRegularSaving(*stateSaveIntervalSec)

		server := NewServer(*listenAddr, NewScheduler(state, *taskTimeoutSec), state)

		fmt.Printf("starting server on %v\n", *listenAddr)

		if err := server.start(*listenAddr); err != nil {
			log.Fatal(err)
		}

	} else {
		fmt.Printf("starting agent\n")

		taskReceiver := func() (PoolTask, error) {

			task, err := NewTaskFromServer(*serverAddr, *agentRequestTimeoutSec)

			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Duration(*agentErrorSleepSec) * time.Second)
				return nil, err
			}

			task.SetJob(func() {
				// sleep random pause
				time.Sleep(time.Duration(rand.Intn(8)) * time.Second)

				solution := NewSolution(task, rand.Float64()*10000, "text representation...")

				// report to server as done
				err := task.ReportToServerWhatDone(*serverAddr, *agentRequestTimeoutSec, solution)
				if err != nil {
					fmt.Println(err)
				}
			})

			return task, nil
		}

		pool := newPool(*agentThreads, taskReceiver)
		pool.Start()
	}
}
