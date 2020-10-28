package utils

import (
	"fmt"
)

type  Worker struct {
	Id int
	Jobs *chan *Task
	Die *chan bool
}

func (w *Worker) Run () { 
	for task := range *w.Jobs {
		if task.Status == STOPPED { // task Cancelled
			fmt.Printf("Worker : %d Not Running Task %d\n", w.Id, task.Id)
			continue
		}

		fmt.Printf("Worker : %d Running Task %d\n", w.Id, task.Id)
		var result []interface{}
		task.SetStatus(RUNNING)
		result = task.Call()
		task.Result = result
		task.SetStatus(FINISHED)
		fmt.Printf("Worker : %d Finished Running Task %d\n", w.Id, task.Id)
	}
}