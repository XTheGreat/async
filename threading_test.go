package async

import (
	"context"
	"testing"
)

func TestThreading(t *testing.T) {
	var acknowledge = []bool{}
	th, err := MakeThreading(Util{
		Read: func(offset, block int) (interface{}, error) { return nil, nil },
		Execute: func(data interface{}, block int) error {
			acknowledge = append(acknowledge, true)
			return nil
		},
		MakeInterval: func() (idx, limit int, err error) { return 1, 200, nil },
		BlockSize:    10,
		Nodes:        5,
	})
	if err != nil {
		t.Error(err)
	}
	th.Go()

	if len(acknowledge) != 20 {
		t.Error("unexpected # of acknowledgements: ", len(acknowledge))
	}
}

func TestThreadingShutdown(t *testing.T) {
	th, err := MakeThreading(Util{
		Read:         func(offset, block int) (interface{}, error) { return []int{}, nil },
		Execute:      func(data interface{}, block int) error { return nil },
		MakeInterval: func() (idx, limit int, err error) { return 0, 100, nil },
		BlockSize:    20,
		Nodes:        5,
	})
	if err != nil {
		t.Error(err)
	}
	th.Go()
	th.Shutdown(context.Background())
}
