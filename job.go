package async

import (
	"context"
	"sync"
)

// Job is an object whose handle can be specified in functions
// to safely manage the execution of multiple goroutines
type Job struct {
	wg sync.WaitGroup
}

// MakeJob safely runs goroutines
func MakeJob() *Job {
	return &Job{}
}

// Run for fitting the graceful Job interface
func (j *Job) Run() {}

// Shutdown waits for running goroutines to return
func (j *Job) Shutdown(ctx context.Context) error {
	j.wg.Wait()
	return nil
}

// Go launches a goroutine and monitors its exit
//
// Example calling a function in a goroutine.
//    j.Go(func(){
//
//        err := SomeFunc(arg1, arg2, ...)
//        if err == nil {
//            // error handling
//        }
//    })
func (j *Job) Go(f func()) {
	j.wg.Add(1)
	go func() {
		f()
		j.wg.Done()
	}()
}
