package utils

import (
	"log"
	"sync"
)

func Parallel(do func(), workers int) {
	wg := &sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			do()
			wg.Done()
		}()
	}
	wg.Wait()
}

func ParallelCheck(do func() error, workers int, errors chan<- error) {
	Parallel(func() {
		err := do()
		if err != nil {
			errors <- err
		}
	}, workers)
}

func ParallelLog(do func() error, workers int) {
	errors := make(chan error)
	go func() {
		for err := range errors {
			log.Printf("parallel error: %s", err)
		}
	}()

	ParallelCheck(do, workers, errors)
	close(errors)
}
