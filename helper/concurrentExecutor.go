package helper

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type TaskFunc func(task interface{}) interface{}

func ConcurrentExecutor(tasks []interface{}, workerFunc TaskFunc, options ...interface{}) []interface{} {
	numWorkers := runtime.NumCPU()
	logTime := false

	for _, opt := range options {
		switch v := opt.(type) {
		case int:
			numWorkers = v
		case bool:
			logTime = v
		}
	}

	taskChan := make(chan interface{}, numWorkers)
	resultChan := make(chan interface{}, len(tasks))
	var wg sync.WaitGroup

	var start time.Time
	if logTime {
		start = time.Now()
	}

	worker := func() {
		for task := range taskChan {
			result := workerFunc(task)
			resultChan <- result
		}
		wg.Done()
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	results := make([]interface{}, 0, len(tasks))
	for res := range resultChan {
		results = append(results, res)
	}

	if logTime {
		fmt.Printf("Total time taken: %v\n", time.Since(start))
	}

	return results
}
