package utils

import (
	"fmt"
	"sync"
	"time"
)

type Future struct {
	ctxt               *Context
	wg                 *sync.WaitGroup
	handleCallbacks    chan bool
	callbacks          []func(*Future)
	callBackExceptions map[int]interface{}
}

func (f *Future) GetId() int {
	return f.ctxt.GetId()
}

func (f *Future) SetContext(ctxt *Context) {
	f.ctxt = ctxt
}

func (f *Future) Cancel() bool {
	return f.ctxt.Cancel()
}

func (f *Future) Cancelled() bool {
	return f.ctxt.GetStatus() == STOPPED
}

func (f *Future) Running() bool {
	return f.ctxt.GetStatus() == RUNNING
}

func (f *Future) Finished() bool {
	return f.ctxt.GetStatus() == FINISHED
}

func (f *Future) recoverFromPanic(callBackId int) {
	if r := recover(); r != nil {
		fmt.Printf("future %d callback %d errored\n", f.GetId(), callBackId)
		f.callBackExceptions[callBackId] = r.(interface{})
	}
}

func (f *Future) AddDoneCallback(callback func(*Future)) {
	defer f.recoverFromPanic(0)

	if f.ctxt.GetStatus() == FINISHED || f.ctxt.GetStatus() == STOPPED {
		callback(f)
	} else {
		f.callbacks = append(f.callbacks, callback)
	}
}

func (f *Future) Exception(timeout ...int) (interface{}, bool) {
	if f.Cancelled() == false || f.Running() == true {
		if timeout == nil {
			<-f.ctxt.Done
			f.ctxt.Done <- true
		} else {
			<-time.After(time.Duration(timeout[0]) * time.Second) // Waiting till the timeout
			if f.Finished() == false || f.Cancelled() == false {
				return nil, false
			}
		}
	}

	return (f.ctxt).GetException(), true
}

func (f *Future) Result(timeout ...int) (interface{}, bool) {
	if f.Cancelled() == false || f.Running() == true {
		if timeout == nil {
			<-f.ctxt.Done
			f.ctxt.Done <- true
		} else {
			<-time.After(time.Duration(timeout[0]) * time.Second) // Waiting till the timeout
		}
	}

	if f.ctxt.GetStatus() == FINISHED {
		return f.ctxt.GetResult(), true
	}

	// No result found: timeout or encountered exception
	return nil, false
}

func (f *Future) callCallback(callBackId int) {
	defer f.recoverFromPanic(callBackId)
	callable := f.callbacks[callBackId]
	callable(f)
}

func (f *Future) invokeCallbacks() {
	<-f.handleCallbacks // Waiting on task to get finished/cancelled

	for callBackId, _ := range f.callbacks {
		f.callCallback(callBackId)
	}

	f.wg.Done()
}

func NewFuture(ftpWaitGroup *sync.WaitGroup) *Future {
	ftr := Future{ctxt: nil,
		wg:                 ftpWaitGroup,
		handleCallbacks:    make(chan bool, 1),
		callbacks:          [](func(*Future)){},
		callBackExceptions: make(map[int]interface{})}
	
	ftr.wg.Add(1)
	go ftr.invokeCallbacks() // Worker to wait on task for running callbacks
	return &ftr
}
