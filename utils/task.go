package utils
import (
	"sync"
)

const (
	READY = 0
	RUNNING = 1
	FINISHED = 2
	STOPPED = 3
)

type Task struct {
	Id int
	Status int // updates should be atomic
	Cancelled bool
	Lock *sync.Mutex
	Call func() (interface{})
	Result interface{}
}

func (t *Task) Cancel () bool {
	t.Lock.Lock()
		if t.Status == READY {
			t.Cancelled = true
			t.Status = STOPPED
			return true
		}
	t.Lock.Unlock()
	return false
}

func (t *Task) SetStatus (status int) {
	t.Lock.Lock()
		t.Status = status
	t.Lock.Unlock()
}

func NewTask (id int, f func()(interface{})) Task {
	return 	Task{ 	Id : id,
					Status : READY,
					Cancelled : false,
					Lock : new(sync.Mutex),
					Call : f,
					Result : nil}
}
