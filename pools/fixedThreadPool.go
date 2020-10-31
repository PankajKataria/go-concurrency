package pools

import (
	"sync"
	u "github.com/PankajKataria/go-concurrency/utils"
)

type ThreadPool interface {
	Submit (func()([]interface{})) (*u.Future)
}

type FixedThreadPool struct {
	JobQueue chan *u.Task
	JobsCount int
	Wg sync.WaitGroup
	DieChan map[int](chan bool)
	ThreadPool // embedding Thread pool interface
}

func (fTP *FixedThreadPool) Submit (callable func() (interface{})) *u.Future {
	ftr := u.NewFuture()
	ctxt := u.NewContext(fTP.JobsCount, &ftr)
	task := u.NewTask(callable, ctxt)
	fTP.JobsCount = fTP.JobsCount + 1
	fTP.JobQueue <- &task // add task to the worker queue
	return &ftr
}

func (fTP *FixedThreadPool) ShutDown () {
	for _, wc := range fTP.DieChan {
		wc <- true // sending dead signal
	}
	fTP.Wg.Wait()
}

func NewFixedThreadPool(workers int) (*FixedThreadPool) {
	cp := FixedThreadPool{	JobQueue : make (chan *u.Task, 100), 
							JobsCount: 0,
							Wg: *new(sync.WaitGroup), 
							DieChan: make(map[int](chan bool))}

	for i := 0; i<workers; i++ {
		cp.Wg.Add(1)
		cchan := make(chan bool, 1)
		cp.DieChan[i] = cchan

		worker := u.Worker{	Id: i, 
							Jobs: &(cp.JobQueue),
							Die: &cchan,
							Wg: &cp.Wg}

		go worker.Run()
	}

	return &cp
}