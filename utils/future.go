package utils

import (
	"time"
)

type Future struct {
	ctxt *Context
	callbacks []func(*Future)()
	callBackExceptions map[int]interface{}
}

func (f *Future) SetContext(ctxt *Context) {
	f.ctxt = ctxt
}

func (f *Future) Cancel () bool {
	return f.ctxt.Cancel()
}

func (f *Future) Cancelled () bool {
	return f.ctxt.GetStatus() == STOPPED
}

func (f *Future) Running () bool {
	return f.ctxt.GetStatus() == RUNNING
}

func (f *Future) Finished () bool {
	return f.ctxt.GetStatus() == FINISHED
}

func (f *Future) AddDoneCallback(callback func(*Future)()) {
	if f.ctxt.GetStatus() == FINISHED || f.ctxt.GetStatus() == STOPPED {
		callback(f)
	} else {
		f.callbacks = append(f.callbacks, callback)
	}
}

func (f *Future) Exception (timeout ...int) interface{} {
	if f.Cancelled () == false || f.Running() == true {
		if (timeout == nil) {
			<-f.ctxt.Done
			f.ctxt.Done<-true
		} else {
			<-time.After(time.Duration(timeout[0]) * time.Second) // Waiting till the timeout
		}
	}
	return (f.ctxt).GetException()
}

func (f *Future) RecoverFromPanic(callBackId int)  {
	if r := recover(); r!=nil {
		f.callBackExceptions[callBackId] = r.(interface{})
	}
}
func (f *Future) InvokeCallBacks () {
	for callBackId, callback := range f.callbacks {
		defer f.RecoverFromPanic(callBackId)
		callback(f);
	}
}

func (f *Future) Result (timeout ...int) (interface{}, bool) {
	if f.Cancelled () == false || f.Running() == true {
		if (timeout == nil) {
			<-f.ctxt.Done
			f.ctxt.Done<-true
		} else {
			<-time.After(time.Duration(timeout[0]) * time.Second) // Waiting till the timeout
		}
	}

	if f.ctxt.GetStatus() == FINISHED {
		return f.ctxt.GetResult(), true
	}

	// Not Finished
	return nil, false
}

func NewFuture() Future {
	return Future{	ctxt: nil,
					callbacks: [](func(*Future)()){},
					callBackExceptions: make(map[int]interface{})}
}
