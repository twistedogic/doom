package cmd

import (
	"os"

	"github.com/twistedogic/doom/cmd/run"
	"github.com/twistedogic/doom/cmd/schedule"
	"github.com/twistedogic/doom/cmd/start"
	"github.com/twistedogic/doom/cmd/token"
	"github.com/urfave/cli"
)

func Start() error {
	app := cli.NewApp()
	app.Name = "doom"
	app.Commands = []cli.Command{
		run.New(),
		start.New(),
		schedule.New(),
		token.New(),
	}
	return app.Run(os.Args)
}
