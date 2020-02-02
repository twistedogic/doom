package start

import (
	"github.com/spf13/afero"
	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/pkg/model/detail"
	"github.com/twistedogic/doom/pkg/model/history"
	"github.com/twistedogic/doom/pkg/model/match"
	"github.com/twistedogic/doom/pkg/model/noop"
	"github.com/twistedogic/doom/pkg/model/odd"
	"github.com/twistedogic/doom/pkg/schedule"
	"github.com/twistedogic/doom/pkg/schedule/job"
	"github.com/twistedogic/doom/pkg/store"
	filestore "github.com/twistedogic/doom/pkg/store/fs"
	"github.com/urfave/cli"
)

var (
	pathFlag string

	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "path, p",
			Value:       ".",
			Usage:       "path for data store",
			Destination: &pathFlag,
		},
	}
)

func Run(c *cli.Context) error {
	fs := afero.NewBasePathFs(afero.NewOsFs(), pathFlag)
	s := store.NewFileStore(filestore.New(fs))
	transformers := []model.TransformFunc{
		odd.Transform,
		history.Transform,
		detail.Transform,
		match.Transform,
		noop.Transform,
	}
	dst := model.New(s, transformers...)
	jobs := make([]schedule.Job, 0, len(configs))
	for _, cfg := range configs {
		jobs = append(jobs, job.New(cfg.Name, cfg.Tap, dst, cfg.Period))
	}
	scheduler := schedule.New(jobs...)
	return scheduler.Start(ContextWithInterrupt())
}

func New() cli.Command {
	return cli.Command{
		Name:   "start",
		Usage:  "start scraping",
		Flags:  flags,
		Action: Run,
	}
}
