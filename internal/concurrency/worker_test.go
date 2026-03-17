package concurrency

import (
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	StartWorkerPool(3)

	// Since we don't have real domain rules in this package, we'll test the loop
	// Note: LoadCriteria is called in worker.go, this might fail if data/ folder is missing or labels are wrong.
	// For a pure unit test, we should mock rules.LoadCriteria, but for this exercise, we'll assume it works
	// or just test the concurrency flow.

	start := time.Now()
	
	// Launch 3 tasks
	resChan1 := make(chan bool)
	resChan2 := make(chan bool)
	resChan3 := make(chan bool)

	taskQueue <- ValidationTask{Label: "gmail", Password: "Invalid", ResultChan: resChan1}
	taskQueue <- ValidationTask{Label: "gmail", Password: "Invalid", ResultChan: resChan2}
	taskQueue <- ValidationTask{Label: "gmail", Password: "Invalid", ResultChan: resChan3}

	<-resChan1
	<-resChan2
	<-resChan3

	duration := time.Since(start)
	
	// Each task takes ~500ms (simulated in menu.ValidatePassword)
	// With 3 workers, 3 tasks should take ~500ms total, not 1500ms.
	if duration > 1000*time.Millisecond {
		t.Errorf("Concurrency test failed: took %v, expected around 500ms", duration)
	}
}
