package main

import (
	"fmt"
	"sync"
	"time"
)

type Pool struct {
	threads                int
	serverAddr             string
	agentName              string
	agentRequestTimeoutSec uint
	agentErrorSleepSec     uint
	submitIntervalSec      uint
	wg                     sync.WaitGroup
	ready                  sync.Mutex
}

func newPool(agentName string, threads uint, serverAddr string, agentRequestTimeoutSec uint, agentErrorSleepSec uint, submitIntervalSec uint) *Pool {
	return &Pool{
		threads:                int(threads),
		serverAddr:             serverAddr,
		agentName:              agentName,
		agentRequestTimeoutSec: agentRequestTimeoutSec,
		agentErrorSleepSec:     agentErrorSleepSec,
		submitIntervalSec:      submitIntervalSec,
	}
}

func (pool *Pool) Start() {

	for i := 0; i < pool.threads; i++ {
		pool.wg.Add(1)
		// start goroutine from pool
		go func(threadNumber int) {
			defer pool.wg.Done()
			fmt.Printf("thread %v started\n", 1+threadNumber)

			solutions := make(chan *Solution)

			for {
				// get next task
				pool.ready.Lock()
				task, err := newTaskFromServer(pool.serverAddr, pool.agentName, pool.agentRequestTimeoutSec)
				pool.ready.Unlock()

				if err != nil {
					fmt.Printf("[%v] newTaskFromServer:%v\n", threadNumber, err)
					time.Sleep(time.Duration(pool.agentErrorSleepSec) * time.Second)
				} else {
					//fmt.Printf("[%v] obtained\n", task.Number)

					task.startWitness(solutions, pool.submitIntervalSec, pool.serverAddr, pool.agentRequestTimeoutSec)

					solutions <- task.do()

					//fmt.Printf("[%v] reported about done task #%v\n", threadNumber, task.Number)
				}

			}
		}(i)
	}

	pool.wg.Wait()
}
