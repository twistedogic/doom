package job

import (
	"context"
	"io"
	"testing"
	"time"
)

type mockTarget struct {
	t *testing.T
}

func NewMockTarget(t *testing.T) mockTarget {
	t.Helper()
	return mockTarget{t}
}

func (m mockTarget) Close() error              { return nil }
func (m mockTarget) Write([]byte) (int, error) { return 0, nil }

type mockTap struct {
	t *testing.T
}

func NewMockTap(t *testing.T) mockTap {
	t.Helper()
	return mockTap{t}
}

func (m mockTap) Update(context.Context, io.WriteCloser) error { return nil }

func TestJob(t *testing.T) {
	src := NewMockTap(t)
	dst := NewMockTarget(t)
	job := New("test", src, dst, time.Hour)
	if err := job.Execute(); err != nil {
		t.Fatal(err)
	}
}
