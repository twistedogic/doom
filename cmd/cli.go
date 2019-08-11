package cmd

import (
	"os"

	"github.com/urfave/cli"
)

func Start() error {
	app := cli.NewApp()
	app.Name = "doom"
	app.Commands = []cli.Command{Run()}
	return app.Run(os.Args)
}
