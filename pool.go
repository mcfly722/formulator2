package main

import (
	"fmt"
	"sync"
)

type PoolTask interface {
	do()
}

type Pool struct {
	threads  int
	wg       sync.WaitGroup
	receiver func() (PoolTask, error)
	ready    sync.Mutex
}

func newPool(threads uint, receiver func() (PoolTask, error)) *Pool {
	return &Pool{
		threads:  int(threads),
		receiver: receiver,
	}
}

func (pool *Pool) Start() {

	for i := 0; i < pool.threads; i++ {
		pool.wg.Add(1)
		// start goroutine from pool
		go func(i int) {
			defer pool.wg.Done()
			fmt.Printf("thread %v started\n", i)

			for {
				// get next task
				pool.ready.Lock()
				task, err := pool.receiver()
				if err != nil {
					fmt.Printf("%v", err)
				}
				pool.ready.Unlock()

				if err == nil {
					fmt.Printf("[%v] started\n", i)
					task.do()
					fmt.Printf("[%v] finished\n", i)
				}
			}

		}(i)
	}

	pool.wg.Wait()
}
