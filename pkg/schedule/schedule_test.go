package schedule

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

type mockOperation struct {
	hasError bool
	called   int
}

func (m *mockOperation) Run(ctx context.Context) error {
	m.called += 1
	if m.hasError {
		return fmt.Errorf("test error")
	}
	return nil
}

func (m *mockOperation) Check(t *testing.T, want int) {
	if m.called != want {
		t.Fatalf("want: %d, got: %d", want, m.called)
	}
}

func TestScheduler(t *testing.T) {
	cases := map[string]struct {
		interval, timeout time.Duration
		want              int
		hasError          bool
	}{
		"base": {
			interval: time.Second,
			timeout:  time.Second,
			want:     5,
			hasError: false,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			mockClock := clock.NewMock()
			m := &mockOperation{tc.hasError, 0}
			s := Scheduler{mockClock, tc.interval, tc.timeout}
			go func() {
				s.Start(ctx, m.Run)
			}()
			for i := 0; i < tc.want; i++ {
				mockClock.Add(tc.interval)
			}
			cancel()
			m.Check(t, tc.want)
		})
	}
}
