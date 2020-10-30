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
	for {
		select { // Attempt at making this non blocking
			case task := <- *((*w).Jobs):
				if task.Status == STOPPED { // task Cancelled
					fmt.Printf("Worker : %d Not Running Task %d\n", w.Id, task.Id)
				} else {
					fmt.Printf("Worker : %d Running Task %d\n", w.Id, task.Id)
					var result interface{}
					task.SetStatus(RUNNING)
					result = task.Call()
					task.Result = result
					task.Done <- true
					task.SetStatus(FINISHED)
					fmt.Printf("Worker : %d Finished Running Task %d\n", w.Id, task.Id)	
				}
			case <- *((*w).Die):
				fmt.Printf("Worker : %d shuting down\n", w.Id)
				w.Wg.Done()
				return
			default:
		}
	}
}