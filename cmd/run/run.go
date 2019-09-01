package run

import (
	"log"

	"github.com/twistedogic/doom/cmd/options"
	"github.com/twistedogic/doom/pkg/config"
	"github.com/urfave/cli"
)

var (
	configFlag string

	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "task config location",
			Destination: &configFlag,
		},
	}
)

func run(tasks []config.Task) error {
	errCh := make(chan error)
	for i := range tasks {
		job, err := options.Load(tasks[i])
		if err != nil {
			return err
		}
		go func() {
			log.Printf("Running %s", job.Name)
			if err := job.Execute(); err != nil {
				errCh <- err
			}
			log.Printf("Complete %s", job.Name)
		}()
	}
	return <-errCh
}

func New() cli.Command {
	run := func(c *cli.Context) error {
		cfg, err := config.Load(configFlag)
		if err != nil {
			return err
		}
		return run(cfg.Tasks)
	}
	return cli.Command{
		Name:   "run",
		Usage:  "run data pipeline",
		Flags:  flags,
		Action: run,
	}
}
