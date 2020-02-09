package job

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/twistedogic/doom/pkg/tap"
)

type mockTarget struct {
	t *testing.T
}

func NewMockTarget(t *testing.T) mockTarget {
	t.Helper()
	return mockTarget{t}
}

func (m mockTarget) Write([]byte) error { return nil }

type mockTap struct {
	t *testing.T
}

func NewMockTap(t *testing.T) mockTap {
	t.Helper()
	return mockTap{t}
}

func (m mockTap) Update(context.Context, tap.Target) error { return nil }

func TestJob_Execute(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	src := NewMockTap(t)
	dst := NewMockTarget(t)
	job := New("test", src, dst, time.Millisecond)
	if err := job.Execute(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestJob_Run(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	src := NewMockTap(t)
	dst := NewMockTarget(t)
	job := New("test", src, dst, time.Millisecond)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	if err := job.Run(ctx); err != nil {
		t.Fatal(err)
	}
}
