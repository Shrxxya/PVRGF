package concurrency

import (
	"context"
	"fmt"
	"sync"
	"PVRGF/internal/logger"
	"PVRGF/internal/menu"
	"PVRGF/internal/rules"
)

var appLog = logger.New("INFO")


// ValidationTask represents a password verification request
type ValidationTask struct {
	Label      string
	Password   string
	ResultChan chan bool
}

var (
	taskQueue   chan ValidationTask
	workerWg    sync.WaitGroup
	workerCtx   context.Context
	cancelFunc  context.CancelFunc
)

// StartWorkerPool starts a specified number of goroutines linked to a context
func StartWorkerPool(numWorkers int, ctx context.Context) {
	workerCtx, cancelFunc = context.WithCancel(ctx)
	appLog.Info("Starting WorkerPool", map[string]interface{}{"workers": numWorkers})
	taskQueue = make(chan ValidationTask, 100)
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go worker(i)
	}
}

// StopWorkerPool gracefully shuts down all background workers
func StopWorkerPool() {
	appLog.Info("Initiating WorkerPool graceful shutdown", nil)
	cancelFunc() // signals all workers to exit
	workerWg.Wait()
	appLog.Info("WorkerPool shut down successfully", nil)
}


func worker(id int) {
	defer workerWg.Done()
	for {
		select {
		case <-workerCtx.Done():
			appLog.Info(fmt.Sprintf("Worker %d shutting down", id), nil)
			return
		case task := <-taskQueue:
			appLog.Info(fmt.Sprintf("Worker %d processing task", id), map[string]interface{}{"label": task.Label})
			criteria, err := rules.LoadCriteria(task.Label)
			valid := false
			if err == nil {
				valid = menu.ValidatePassword(task.Password, criteria)
			}
			appLog.Info(fmt.Sprintf("Worker %d finished task", id), map[string]interface{}{
				"label": task.Label,
				"valid": valid,
			})
			task.ResultChan <- valid
		}
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
