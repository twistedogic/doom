package token

import (
	"github.com/twistedogic/doom/pkg/target/drive/token"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

var (
	credFlag  string
	cacheFlag string

	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "cred, c",
			Value:       "./credential.json",
			Usage:       "google credential",
			Destination: &credFlag,
		},
		cli.StringFlag{
			Name:        "cache, o",
			Value:       "./cache.json",
			Usage:       "token cache",
			Destination: &cacheFlag,
		},
	}
)

func New() cli.Command {
	run := func(c *cli.Context) error {
		ctx := context.Background()
		scope := drive.DriveScope
		_, err := token.GetClient(ctx, scope, credFlag, cacheFlag)
		return err
	}
	return cli.Command{
		Name:   "token",
		Usage:  "get google oauth token",
		Flags:  flags,
		Action: run,
	}
}
