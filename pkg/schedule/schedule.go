package schedule

import (
	"context"
	"log"
	"time"

	"github.com/benbjohnson/clock"
)

type Scheduler struct {
	clock.Clock
	Interval, Timeout time.Duration
}

type Task func(context.Context) error

func New(interval, timeout time.Duration) Scheduler {
	return Scheduler{
		Clock:    clock.New(),
		Interval: interval,
		Timeout:  timeout,
	}
}

func (s Scheduler) runOnce(ctx context.Context, task Task) {
	execCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-s.After(s.Timeout)
		cancel()
	}()
	if err := task(execCtx); err != nil {
		log.Println(err)
	}
}

func (s Scheduler) Start(ctx context.Context, task Task) error {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			go s.runOnce(ctx, task)
		}
		s.Sleep(s.Interval)
	}
	return ctx.Err()
}
