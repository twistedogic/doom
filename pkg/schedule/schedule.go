package schedule

import (
	"context"
	"log"
	"sync"
)

type Job interface {
	Run(context.Context) error
}

type Scheduler struct {
	entries []Job
}

func New(entries ...Job) Scheduler {
	return Scheduler{entries}
}

func (s Scheduler) Start(ctx context.Context) error {
	wg := new(sync.WaitGroup)
	errCh := make(chan error)
	defer close(errCh)
	go func() {
		for err := range errCh {
			log.Println(err)
		}
	}()
	for _, job := range s.entries {
		jobCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		wg.Add(1)
		go func(c context.Context, j Job) {
			defer wg.Done()
			if err := j.Run(c); err != nil {
				errCh <- err
			}
		}(jobCtx, job)
	}
	wg.Wait()
	return nil
}
