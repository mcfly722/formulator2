package main

import (
	"fmt"
	"sync"
)

type PoolTask interface {
	Do()
}

type Pool struct {
	queue    chan PoolTask
	capacity int
	receiver func() (PoolTask, error)
	wg       sync.WaitGroup
}

func newPool(capacity uint, receiver func() (PoolTask, error)) *Pool {
	return &Pool{
		queue:    make(chan PoolTask, capacity),
		capacity: int(capacity),
		receiver: receiver,
	}
}

func (pool *Pool) Start() {
	go func() {
		for {
			// take new task from server
			obtainedTask, err := pool.receiver()
			if err == nil {
				pool.queue <- obtainedTask
			}
		}
	}()

	for i := 0; i < pool.capacity; i++ {
		pool.wg.Add(1)
		// start goroutine from pool
		go func(i int) {

			defer pool.wg.Done()

			fmt.Printf("thread %v started\n", i)
			for {
				// get next task from queue
				task := <-pool.queue
				fmt.Printf("thread %v taked task from queue\n", i)
				task.Do()
				fmt.Printf("thread %v finished task\n", i)
			}

		}(i)
	}

	pool.wg.Wait()
}
