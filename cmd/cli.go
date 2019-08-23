package cmd

import (
	"os"

	"github.com/twistedogic/doom/cmd/run"
	"github.com/twistedogic/doom/cmd/start"
	"github.com/urfave/cli"
)

func Start() error {
	app := cli.NewApp()
	app.Name = "doom"
	app.Commands = []cli.Command{run.Run(), start.Run()}
	return app.Run(os.Args)
}
