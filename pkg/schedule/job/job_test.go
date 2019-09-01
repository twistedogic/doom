package job

import (
	"testing"
	"time"

	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/target"
)

type mockTarget struct {
	t *testing.T
}

func NewMockTarget(t *testing.T) mockTarget {
	t.Helper()
	return mockTarget{t}
}

func (m mockTarget) Load(config.Setting) error    { return nil }
func (m mockTarget) UpsertItem(interface{}) error { return nil }
func (m mockTarget) BulkUpsert(interface{}) error { return nil }
func (m mockTarget) GetLastUpdate() time.Time     { return time.Now() }

type mockTap struct {
	t *testing.T
}

func NewMockTap(t *testing.T) mockTap {
	t.Helper()
	return mockTap{t}
}

func (m mockTap) Load(config.Setting) error  { return nil }
func (m mockTap) Update(target.Target) error { return nil }

func TestJob(t *testing.T) {
	src := NewMockTap(t)
	dst := NewMockTarget(t)
	job := New("test", src, dst, time.Hour)
	if err := job.Execute(); err != nil {
		t.Fatal(err)
	}
}
