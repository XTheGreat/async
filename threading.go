package async

import (
	"context"
	"errors"
	"sync"

	"github.com/golang/glog"
)

// Threading is an organized collection of goroutines that
// executes a single task concurrently.
type Threading struct {
	mu         sync.Mutex
	pos, limit int
	sigterm    chan (int)

	// Read retrieves data in blocks
	Read func(offset, block int) (interface{}, error)
	// Execute processes data
	Execute func(data interface{}, block int) error
	// BlockSize of data processed in a single round of execution.
	BlockSize int
	// Nodes refers to number of running goroutines.
	Nodes int
}

type Util struct {
	// Read retrieves data in blocks
	Read func(offset, block int) (interface{}, error)
	// Execute processes data
	Execute func(data interface{}, block int) error
	// MakeInterval returns the offset and limit [a,b]
	// of data to be processed
	MakeInterval func() (idx, limit int, err error)
	// BlockSize of data processed in a single round of execution.
	BlockSize int
	// Nodes refers to number of running goroutines
	Nodes int
}

// MakeThreading safely executes multiple goroutines to process a single task!
// It includes a Read and Execute operation, and expects offset & limit index of
// data to be processed.
func MakeThreading(d Util) (*Threading, error) {
	if d.Read == nil || d.Execute == nil ||
		d.MakeInterval == nil || d.BlockSize == 0 || d.Nodes == 0 {
		return nil, errors.New("invalid arguments")
	}

	idx, limit, err := d.MakeInterval()
	if err != nil {
		return nil, err
	}
	return &Threading{
		sigterm: make(chan int, 1),
		pos:     idx,
		limit:   limit,

		Read:      d.Read,
		Execute:   d.Execute,
		BlockSize: d.BlockSize,
		Nodes:     d.Nodes,
	}, nil
}

// Run for fitting the graceful Job interface
func (th *Threading) Run() {}

// Go executes tasks in multiple go-routines
func (th *Threading) Go() {
	var wg sync.WaitGroup
	glog.Info("Main: Running task")
	for i := 0; i < th.Nodes; i++ {
		wg.Add(1)
		go th.work(&wg, i)
	}

	wg.Wait()
	glog.Info("Main: Successfully finished task!")
}

func (th *Threading) work(wg *sync.WaitGroup, idx int) error {
	defer wg.Done()

	glog.Infof("Thread running: #%d", idx)
	err := th.do()
	if err != nil {
		glog.Error(err)
	}

	glog.Infof("Thread closed: #%d", idx)
	return nil
}

func (th *Threading) done() bool {
	select {
	case <-th.sigterm:
		return true
	default:
		return th.pos >= th.limit
	}
}

func (th *Threading) do() error {
	for !th.done() {

		/** critical section entry, avoid race!! */
		th.mu.Lock()

		offset := th.pos
		th.pos += th.BlockSize

		th.mu.Unlock()
		/** critical section exit */

		data, err := th.Read(offset, th.BlockSize)
		if err != nil {
			return err
		}

		err = th.Execute(data, th.BlockSize)
		if err != nil {
			return err
		}
	}
	return nil
}

// Shutdown waits for running goroutines to return
func (th *Threading) Shutdown(ctx context.Context) error {
	glog.Info("Main: Sigterm received!")
	close(th.sigterm)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
