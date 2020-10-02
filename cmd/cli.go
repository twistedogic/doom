package cmd

import (
	"os"

	"github.com/twistedogic/doom/cmd/token"
	"github.com/twistedogic/doom/cmd/update"
	"github.com/urfave/cli/v2"
)

func Start() error {
	app := cli.NewApp()
	app.Name = "doom"
	app.Commands = []*cli.Command{
		token.New(),
		update.New(),
	}
	return app.Run(os.Args)
}
