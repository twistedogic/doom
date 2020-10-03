package result

import (
	"time"

	"github.com/twistedogic/doom/cmd/update/updateutil"
	"github.com/twistedogic/doom/pkg/tap/jockey"
	"github.com/urfave/cli/v2"
)

const (
	dateFormat = "2006-01-02"
	month      = 30 * 24 * time.Hour
)

var (
	startFlag, endFlag string

	resultFlags = []cli.Flag{
		&cli.StringFlag{
			Name:        "start",
			Usage:       "start date (YYYY-MM-DD)",
			Value:       time.Now().Add(-month).Format(dateFormat),
			Destination: &startFlag,
		},
		&cli.StringFlag{
			Name:        "end",
			Usage:       "end date (YYYY-MM-DD)",
			Value:       time.Now().Format(dateFormat),
			Destination: &endFlag,
		},
	}
)

func run(c *cli.Context) error {
	store, err := updateutil.Store()
	if err != nil {
		return err
	}
	start, err := time.Parse(dateFormat, startFlag)
	if err != nil {
		return err
	}
	end, err := time.Parse(dateFormat, endFlag)
	if err != nil {
		return err
	}
	tap := jockey.NewResultTap(updateutil.BaseURLFlag, updateutil.RateFlag, start, end)
	return tap.Update(c.Context, store)
}

func New() *cli.Command {
	return &cli.Command{
		Name:   "result",
		Usage:  "update result record",
		Flags:  append(updateutil.Flags, resultFlags...),
		Action: run,
	}
}
