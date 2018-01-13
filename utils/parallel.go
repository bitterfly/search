package utils

import (
	"log"
	"sync"
)

func Parallel(work func(), workers int) {
	wg := &sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			work()
			wg.Done()
		}()
	}
	wg.Wait()
}

func ParallelCheck(work func() error, workers int, errors chan<- error) {
	Parallel(func() {
		err := work()
		if err != nil {
			errors <- err
		}
	}, workers)
}

func ParallelLog(work func() error, workers int) {
	errors := make(chan error)
	go func() {
		for err := range errors {
			log.Printf("parallel error: %s", err)
		}
	}()

	ParallelCheck(work, workers, errors)
	close(errors)
}
