package async

import (
	"testing"
)

func TestBackground(t *testing.T) {
	background := NewListener()

	job := MakeJob()
	background.Observe(job)

	thread, err := MakeThreading(Util{
		Read: func(offset, block int) (interface{}, error) { return nil, nil },
		Execute: func(data interface{}, block int) error {
			for i := 0; i < block; i++ {
				// do something
			}
			return nil
		},
		MakeInterval: func() (idx, limit int, err error) { return 1, 2000, nil },
		BlockSize:    10,
		Nodes:        5,
	})
	if err != nil {
		t.Error(err)
	}
	background.Observe(thread)
	thread.Go()

	background.Quit()

	c := make(chan interface{}, 1)
	job.Go(func() {
		close(c)
	})
}
