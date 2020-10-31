package utils

import (
	// "fmt"
	"sync"
)

const (
	READY    = 0
	RUNNING  = 1
	FINISHED = 2
	STOPPED  = 3
)

type Context struct {
	id        int
	status    int // updates should be atomic
	Cancelled bool
	Lock      *sync.Mutex
	Done      chan bool
	exception interface{}
	result    interface{}
	future    *Future
}

func (ctxt *Context) GetId() int {
	return ctxt.id
}

func (ctxt *Context) GetStatus() int {
	return ctxt.status
}

func (ctxt *Context) GetResult() interface{} {
	return ctxt.result
}

func (ctxt *Context) GetException() interface{} {
	return ctxt.exception
}

func (ctxt *Context) SetStatus(status int) {
	ctxt.Lock.Lock()
	ctxt.status = status
	ctxt.Lock.Unlock()
}

func (ctxt *Context) SetResult(result interface{}) {
	ctxt.Lock.Lock()
	ctxt.result = result
	ctxt.Lock.Unlock()
}

func (ctxt *Context) SetException(exception interface{}) {
	ctxt.Lock.Lock()
	ctxt.exception = exception
	ctxt.Lock.Unlock()
}

func (cntxt *Context) Cancel() bool {
	cntxt.Lock.Lock()
	if cntxt.status == READY {
		cntxt.Cancelled = true
		cntxt.status = STOPPED
		return true
	}
	cntxt.Lock.Unlock()
	return false
}

func NewContext(id int, f *Future) *Context {
	ctxt := Context{id: id,
		status:    READY,
		Cancelled: false,
		Lock:      new(sync.Mutex),
		Done:      make(chan bool, 1),
		future:    f}

	(*f).SetContext(&ctxt)
	return &ctxt
}
