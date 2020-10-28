package pools

import (
	f "github.com/PankajKataria/go-concurrency/future"
	u "github.com/PankajKataria/go-concurrency/utils"
)

type ThreadPool interface {
	Submit (func()([]interface{})) (*f.Future)
}

type FixedThreadPool struct {
	JobQueue chan *u.Task
	jobsDone int
	ThreadPool // embedding Thread pool interface
}

func (fTP *FixedThreadPool) Submit (callable func() ([]interface{})) *f.Future {
	nt := u.NewTask(fTP.jobsDone, callable) // current Task Object
	fTP.JobQueue <- &nt // add task to the worker queue
	return &f.Future{Task:&nt} // current future Object
} 

func NewFixedThreadPool(workers int) (*FixedThreadPool) {
	cp := FixedThreadPool{JobQueue : make (chan *u.Task, 100)}
	for i := 0; i<workers; i++ {
		worker := u.Worker{Id: i, Jobs: &cp.JobQueue, Die: new(chan bool)}
		go worker.Run()
	}
	return &cp
}