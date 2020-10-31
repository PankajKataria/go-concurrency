package utils

import (
	"fmt"
)

type Task struct {
	cntxt * Context
	Call func() (interface{})
}

func RecoverFromPanic (t *Task, workerId, tId int) {
	if r := recover(); r != nil {
		fmt.Printf("Worker %d: Task  %d Paniced \n", workerId, tId)
		(t.cntxt).SetException(r)
		t.cntxt.SetStatus(STOPPED)
		t.cntxt.Done <- true
	}
}

func (t *Task) Execute (workerId int) {
	defer RecoverFromPanic(t, workerId, t.GetId())

	if t.cntxt.GetStatus() != STOPPED {
		t.cntxt.SetStatus(RUNNING)
		fmt.Printf("Worker %d: Running Task %d\n", workerId, t.GetId())

		callable := t.Call
		result := callable()

		t.cntxt.SetResult(result)
		t.cntxt.SetStatus(FINISHED)
		t.cntxt.Done <- true // signal waiting consumers for Result/Exception

		t.ExecuteCallbacks() // execute Callbacks
		fmt.Printf("Worker %d: Finished Running Task %d\n", workerId, t.GetId())
	} else {
		fmt.Printf("Worker %d: Task %d is Cancelled\n", workerId, t.GetId())
	}
}

func (t *Task) ExecuteCallbacks () {
	t.cntxt.future.InvokeCallBacks()
}

func (t *Task) GetId() int {
	return t.cntxt.GetId()
}

func NewTask (f func()(interface{}), cntxt *Context) Task {
	return 	Task{	cntxt: cntxt,
					Call : f}
}
