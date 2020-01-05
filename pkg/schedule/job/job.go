package job

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/twistedogic/doom/pkg/tap"
)

type Status string

const (
	ERROR Status = "ERROR"
	DONE  Status = "DONE"
)

type Job struct {
	Name     string
	Src      tap.Tap
	Dst      io.WriteCloser
	Interval time.Duration
}

func New(name string, src tap.Tap, dst io.WriteCloser, interval time.Duration) *Job {
	return &Job{
		Name:     name,
		Src:      src,
		Dst:      dst,
		Interval: interval,
	}
}

func (j Job) Execute(ctx context.Context) error {
	return j.Src.Update(ctx, j.Dst)
}

func (j Job) Run(ctx context.Context) error {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			runCtx, cancel := context.WithTimeout(ctx, j.Interval)
			defer cancel()
			if err := j.Execute(runCtx); err != nil {
				log.Printf("%s %s %v", ERROR, j.Name, err)
			} else {
				log.Printf("%s %s", DONE, j.Name)
			}
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}
