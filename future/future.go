package future

import (
	"time"
	u "github.com/PankajKataria/go-concurrency/utils"
)

type Future struct {
	Task *u.Task
}

func (f *Future) Cancel () bool {
	return f.Task.Cancel()
}

func (f *Future) Cancelled () bool {
	return f.Task.Cancelled
}

func (f *Future) Running () bool {
	return f.Task.Status == u.RUNNING; 
}

func (f *Future) Finished () bool {
	return f.Task.Status == u.FINISHED; 
}

func (f *Future) Result (timeout int) (interface{}, bool) {
	if f.Cancelled () == false && f.Running() == true {
		<-time.After(time.Duration(timeout) * time.Second) // Waiting till the timeout
	}

	if f.Task.Status == u.FINISHED {
		return f.Task.Result, true
	}

	return nil, false
}


