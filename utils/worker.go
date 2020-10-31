package utils

import (
	"fmt"
	"sync"
)

type  Worker struct {
	Id int
	Jobs *chan *Task
	Die *chan bool
	Wg *sync.WaitGroup
}

func (w *Worker) Run () {
	fmt.Printf("Worker %d: Alive\n", w.Id)
	for {
		select { // Attempt at making this non blocking
			case task := <- *(w.Jobs):
					task.Execute(w.Id)
			case <- *((*w).Die):
				w.Wg.Done()
				fmt.Printf("Worker %d: Shuting down\n", w.Id)
				return
			default:
		}
	}
}