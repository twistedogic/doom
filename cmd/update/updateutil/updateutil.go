package updateutil

import (
	store "github.com/twistedogic/doom/pkg/store/badger"
	tap "github.com/twistedogic/doom/pkg/tap/jockey"
	"github.com/twistedogic/doom/pkg/transform"
	"github.com/twistedogic/doom/pkg/transform/jockey"
	"github.com/urfave/cli/v2"
)

var (
	PathFlag, BaseURLFlag string
	RateFlag              int

	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "path",
			Value:       ".",
			Usage:       "path for data store",
			Destination: &PathFlag,
		},
		&cli.StringFlag{
			Name:        "base",
			Value:       tap.BaseURL,
			Usage:       "base url",
			Destination: &BaseURLFlag,
		},
		&cli.IntFlag{
			Name:        "rate",
			Value:       5,
			Usage:       "request rate (per second)",
			Destination: &RateFlag,
		},
	}
)

func Store() (transform.Transformer, error) {
	s, err := store.New(PathFlag)
	return transform.New(jockey.TransformJockey, s), err
}
