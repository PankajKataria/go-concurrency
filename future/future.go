package future

import (
	u "github.com/PankajKataria/go-concurrency/utils"
	"time"
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

func (f *Future) Result (timeout int) ([]interface{}, bool) {
	if f.Cancelled() == true {
		return nil, false
	}

	if f.Cancelled () == false && f.Running() == true {
		<-time.After(time.Duration(timeout) * time.Second)
	}

	if f.Task.Status == u.FINISHED {
		return f.Task.Result, true
	}

	return nil, false
}


