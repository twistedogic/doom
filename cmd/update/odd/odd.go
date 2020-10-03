package odd

import (
	"github.com/twistedogic/doom/cmd/update/updateutil"
	"github.com/twistedogic/doom/pkg/tap/jockey"
	"github.com/urfave/cli/v2"
)

var (
	betFlag string

	oddFlags = []cli.Flag{
		&cli.StringFlag{
			Name:        "bet",
			Value:       "had",
			Usage:       "bet type",
			Destination: &betFlag,
		},
	}
)

func run(c *cli.Context) error {
	store, err := updateutil.Store()
	if err != nil {
		return err
	}
	tap := jockey.NewOddTap(updateutil.BaseURLFlag, betFlag, updateutil.RateFlag)
	return tap.Update(c.Context, store)
}

func New() *cli.Command {
	return &cli.Command{
		Name:   "odd",
		Usage:  "update odd record",
		Flags:  append(updateutil.Flags, oddFlags...),
		Action: run,
	}
}
