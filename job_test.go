package async

import (
	"context"
	"testing"
)

func BenchmarkJob(b *testing.B) {
	dummyFunc := func() {}
	job := MakeJob()
	b.ResetTimer()
	// run the amount we want to run
	for n := 0; n < b.N; n++ {
		// we're sending a dummy job to hopefully just capture the amount of time it takes to execute a func and
		// not the time a func would take to run since that's the user's responsibility
		job.Go(dummyFunc)
	}
	// also wait for it to shutdown
	job.Shutdown(context.Background())
}

func TestJob(t *testing.T) {
	c := make(chan interface{}, 1)
	job := MakeJob()
	job.Go(func() {
		close(c)
	})
	<-c
	job.Shutdown(context.Background())
}
