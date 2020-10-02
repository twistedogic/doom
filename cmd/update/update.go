package update

import (
	"github.com/twistedogic/doom/cmd/update/odd"
	"github.com/twistedogic/doom/cmd/update/result"
	"github.com/urfave/cli/v2"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "update match records",
		Subcommands: []*cli.Command{
			odd.New(),
			result.New(),
		},
	}
}
