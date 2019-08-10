package schedule

import (
	"log"

	"github.com/robfig/cron"
	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/pkg/target"
)

func ToFunc(src tap.Tap, dst target.Target) func() {
	return func() {
		if err := src.Update(dst); err != nil {
			log.Println(err)
			return
		}
		log.Println("Done")
	}
}

type Scheduler struct {
	*cron.Cron
}

func New() *Scheduler {
	return &Scheduler{cron.New()}
}

func (s *Scheduler) Add(spec string, src tap.Tap, dst target.Target) error {
	return s.AddFunc(spec, ToFunc(src, dst))
}
