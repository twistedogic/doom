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
	Dst     io.WriteCloser
}

func New(name string, src tap.Tap, dst io.WriteCloser, timeout time.Duration) *Job {
	return &Job{
		Name:    name,
		Src:     src,
		Dst:     dst,
		Timeout: timeout,
	}
}

func (j *Job) Set(src tap.Tap, dst io.WriteCloser) {
	j.SetSrc(src)
	j.SetDst(dst)
}

func (j *Job) SetSrc(src tap.Tap) {
	j.Src = src
}

func (j *Job) SetDst(dst io.WriteCloser) {
	j.Dst = dst
}

func (j *Job) Execute() error {
	ctx, _ := context.WithTimeout(context.Background(), j.Timeout)
	return j.Src.Update(ctx, j.Dst)
}

func (j *Job) Run() {
	timeit := TimeIt()
	if err := j.Execute(); err != nil {
		log.Printf("%s %s %s %v", ERROR, j.Name, timeit(), err)
		return
	}
	log.Printf("%s %s %s", DONE, j.Name, timeit())
}
