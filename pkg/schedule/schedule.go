package schedule

import (
	"github.com/robfig/cron"
	"github.com/twistedogic/doom/pkg/schedule/job"
)

type Scheduler struct {
	*cron.Cron
}

func New() *Scheduler {
	return &Scheduler{cron.New()}
}

func (s *Scheduler) Add(spec string, j *job.Job) error {
	return s.AddFunc(spec, j.ToFunc())
}
