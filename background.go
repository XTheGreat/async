package async

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

type Task interface {
	Shutdown(ctx context.Context) error
}

type Background struct {
	tasks []Task
	quit  bool
}

func NewListener() *Background {
	b := &Background{}
	b.listen()
	return b
}

func (b *Background) Observe(task Task) {
	b.tasks = append(b.tasks, task)
}

func (b *Background) Quit() {
	if !b.quit {
		b.shutdown()
	}
}

func (b *Background) shutdown() {
	b.quit = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	glog.Info("Shutting down background jobs (^C to exit now)")
	for _, task := range b.tasks {
		err := task.Shutdown(ctx)
		if err != nil {
			glog.Error(err)
		}
	}
}

func (b *Background) listen() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	var stop bool
	go func() {
		for !stop {
			s := <-quit
			switch s {
			case os.Interrupt, syscall.SIGTERM:
				stop = true
				b.shutdown()
			}
		}
	}()
}
