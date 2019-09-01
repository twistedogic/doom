package job

import (
	"context"
	"log"
	"time"

	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/pkg/target"
)

type Status string

const (
	ERROR   Status = "ERROR"
	TIMEOUT Status = "TIMEOUT"
	DONE    Status = "DONE"
)

func TimeIt() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}

type Job struct {
	Timeout time.Duration
	Name    string
	Src     tap.Tap
	Dst     target.Target
}

func New(name string, src tap.Tap, dst target.Target, timeout time.Duration) *Job {
	return &Job{
		Name:    name,
		Src:     src,
		Dst:     dst,
		Timeout: timeout,
	}
}

func (j *Job) Set(src tap.Tap, dst target.Target) {
	j.SetSrc(src)
	j.SetDst(dst)
}

func (j *Job) SetSrc(src tap.Tap) {
	j.Src = src
}

func (j *Job) SetDst(dst target.Target) {
	j.Dst = dst
}

func (j *Job) Execute() error {
	return j.Src.Update(j.Dst)
}

func (j *Job) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
	defer cancel()
	timeit := TimeIt()
	errCh := make(chan error)
	go func() {
		errCh <- j.Execute()
	}()
	select {
	case <-ctx.Done():
		log.Printf("%s %s %s", TIMEOUT, j.Name, timeit())
		return
	case err := <-errCh:
		if err != nil {
			log.Printf("%s %s %s %v", ERROR, j.Name, timeit(), err)
			return
		}
		log.Printf("%s %s %s", DONE, j.Name, timeit())
	}
}
