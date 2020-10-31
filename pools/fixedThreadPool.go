package pools

import (
	"fmt"
	"sync"
	u "github.com/PankajKataria/go-concurrency/utils"
)

type ThreadPool interface {
	Submit (func()([]interface{})) (*u.Future, error)
}

type FixedThreadPool struct {
	dead bool
	JobQueue chan *u.Task
	JobsCount int
	workersSync sync.WaitGroup
	taskSync sync.WaitGroup
	DieChan map[int](chan bool)
	ThreadPool // embedding Thread pool interface
}

func (fTP *FixedThreadPool) Submit (callable func() (interface{})) (*u.Future, error) {
	if (fTP.dead == true) {
		return nil, fmt.Errorf("Thread Pool is Dead")
	}

	ftr := u.NewFuture(&fTP.taskSync)
	ctxt := u.NewContext(fTP.JobsCount, ftr)
	task := u.NewTask(callable, ctxt)
	fTP.JobsCount = fTP.JobsCount + 1
	fTP.JobQueue <- &task // add task to the worker queue
	// fmt.Printf ("Submit %p : %v\n", ftr, *ftr)
	return ftr, nil
}

func (fTP *FixedThreadPool) ShutDown () {
	fTP.dead = true

	for _, wc := range fTP.DieChan {
		wc <- true // sending dead signal
	}

	fTP.workersSync.Wait()	// Wait till all workers are down
	fTP.taskSync.Wait() // Wait till all tasks + callbacks are processed
}

func NewFixedThreadPool(workers int) (*FixedThreadPool, error) {
	if (workers < 1) {
		return nil, fmt.Errorf("Workers > 1 Required")
	}

	cp := FixedThreadPool{	JobQueue : make (chan *u.Task, 100), 
							JobsCount: 0,
							workersSync: *new(sync.WaitGroup),
							taskSync: *new(sync.WaitGroup),
							DieChan: make(map[int](chan bool))}

	cp.workersSync.Add(workers) // wait Group On Task Added
	for i := 0; i<workers; i++ {
		cchan := make(chan bool, 1)
		cp.DieChan[i] = cchan
		worker := u.Worker{	Id: i, 
							Jobs: &(cp.JobQueue),
							Die: &cchan,
							Wg: &(cp.workersSync)}

		go worker.Run()
	}

	return &cp, nil
}