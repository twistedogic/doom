package schedule

import (
	"fmt"
	"log"
	"net/http"

	"github.com/twistedogic/doom/cmd/options"
	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/schedule"
	"github.com/urfave/cli"
)

var (
	configFlag string
	portFlag   int

	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "task config location",
			Destination: &configFlag,
		},
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "port for report service",
			Value:       3000,
			Destination: &portFlag,
		},
	}
)

func scheduleTask(s *schedule.Scheduler, tasks []config.Task) error {
	for _, t := range tasks {
		job, err := options.Load(t)
		if err != nil {
			return err
		}
		log.Printf("Scheduling %s", job.Name)
		if err := s.AddJob(t.Schedule, job); err != nil {
			return err
		}
	}
	return nil
}

func New() cli.Command {
	run := func(c *cli.Context) error {
		cfg, err := config.Load(configFlag)
		if err != nil {
			return err
		}
		scheduler := schedule.New()
		if err := scheduleTask(scheduler, cfg.Tasks); err != nil {
			return err
		}
		log.Printf("scheduled %d tasks", len(cfg.Tasks))
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", portFlag),
			Handler: scheduler,
		}
		scheduler.Start()
		log.Printf("Server Running at %d", portFlag)
		if err := server.ListenAndServe(); err != nil {
			scheduler.Stop()
			return err
		}
		return nil
	}
	return cli.Command{
		Name:   "schedule",
		Usage:  "schedule and run cron jobs",
		Flags:  flags,
		Action: run,
	}
}
