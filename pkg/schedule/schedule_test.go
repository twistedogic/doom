package schedule

import (
	"context"
	"testing"
	"time"
)

type mockJob struct{}

func (m mockJob) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func setupJobs(n int) []Job {
	jobs := make([]Job, n)
	for i := range jobs {
		jobs[i] = mockJob{}
	}
	return jobs
}

func TestScheduler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	jobs := setupJobs(5)
	s := New(jobs...)
	if err := s.Start(ctx); err != nil {
		t.Fatal(err)
	}
}
