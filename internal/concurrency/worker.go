package concurrency

import (
	"fmt"
	"PVRGF/internal/menu"
	"PVRGF/internal/rules"
)

// ValidationTask represents a password verification request
type ValidationTask struct {
	Label      string
	Password   string
	ResultChan chan bool
}

var taskQueue chan ValidationTask

// StartWorkerPool starts a specified number of goroutines to process validation tasks
func StartWorkerPool(numWorkers int) {
	fmt.Printf("[WorkerPool] Starting %d workers...\n", numWorkers)
	taskQueue = make(chan ValidationTask, 100)
	for i := 0; i < numWorkers; i++ {
		go worker(i)
	}
}

func worker(id int) {
	for task := range taskQueue {
		fmt.Printf("[Worker %d] Received task for label: %s\n", id, task.Label)
		criteria, err := rules.LoadCriteria(task.Label)
		valid := false
		if err == nil {
			valid = menu.ValidatePassword(task.Password, criteria)
		}
		fmt.Printf("[Worker %d] Finished task for label: %s (Valid: %v)\n", id, task.Label, valid)
		task.ResultChan <- valid
	}
}

// SubmitTask sends a validation task to the worker pool and returns the result
func SubmitTask(label, password string) bool {
	resultChan := make(chan bool)
	taskQueue <- ValidationTask{
		Label:      label,
		Password:   password,
		ResultChan: resultChan,
	}
	return <-resultChan
}
