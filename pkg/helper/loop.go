package helper

import (
	"sync"
)

type Loop struct {
	wg      *sync.WaitGroup
	errCh   chan error
	doneCh  chan struct{}
	running chan struct{}
	limit   int
}

func NewLoop(limit int) *Loop {
	errCh := make(chan error)
	doneCh := make(chan struct{})
	wg := new(sync.WaitGroup)
	running := make(chan struct{}, 0)
	if limit > -1 {
		running = make(chan struct{}, limit)
	}
	return &Loop{wg, errCh, doneCh, running, limit}
}

func (l *Loop) Add(n int) {
	if l.limit != -1 {
		for i := 0; i < n; i++ {
			l.running <- struct{}{}
		}
	}
	l.wg.Add(n)
}

func (l *Loop) Done(err error) {
	if l.limit != -1 {
		<-l.running
	}
	defer l.wg.Done()
	if err != nil {
		l.errCh <- err
	}
}

func (l *Loop) Wait() error {
	go func() {
		l.wg.Wait()
		l.doneCh <- struct{}{}
	}()
	select {
	case err := <-l.errCh:
		return err
	case <-l.doneCh:
		return nil
	}
}
