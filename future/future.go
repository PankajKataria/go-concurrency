package future

import (
	"sync"
	"time"
	u "github.com/PankajKataria/go-concurrency/utils"
)

type Future struct {
	Task *u.Task
	Callbacks []func(*Future)()
	FException interface{}
	CWG *sync.WaitGroup
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

func (f *Future) RecoverFromException () {
	if r := recover(); r!=nil {
		f.FException = r.(interface{})
	}
	f.InvokeCallBacks()
}

func (f *Future) AddDoneCallback(callback func(*Future)()) {
	if f.Task.Status == u.FINISHED || f.Task.Status == u.STOPPED {
		callback(f)
	} else {
		f.Callbacks = append(f.Callbacks, callback)
	}
}

func (f *Future) Exception (timeout interface{}) interface{} {
	defer f.RecoverFromException()
	if f.Cancelled () == false || f.Running() == true {
		if (timeout == nil) {
			<-f.Task.Done
		} else {
			<-time.After(time.Duration(timeout.(int)) * time.Second) // Waiting till the timeout
		}
	}

	return f.FException
}

func (f *Future) InvokeCallBacks () {
	for _, callback := range f.Callbacks {
		callback(f);
	}
}

func (f *Future) Result (timeout interface{}) (interface{}, bool) {
	defer f.RecoverFromException()

	if f.Cancelled () == false || f.Running() == true {
		if (timeout == nil) {
			<-f.Task.Done
		} else {
			<-time.After(time.Duration(timeout.(int)) * time.Second) // Waiting till the timeout
		}
	}

	if f.Task.Status == u.FINISHED {
		return f.Task.Result, true
	}

	// Not Finished
	return nil, false
}


