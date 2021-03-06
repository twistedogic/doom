package token

import (
	"github.com/twistedogic/doom/pkg/store/fs/drive"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
	gdrive "google.golang.org/api/drive/v3"
)

var (
	credFlag  string
	cacheFlag string

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "cred",
			Value:       "./credential.json",
			Usage:       "google credential",
			Destination: &credFlag,
		},
		&cli.StringFlag{
			Name:        "cache",
			Value:       "./cache.json",
			Usage:       "token cache",
			Destination: &cacheFlag,
		},
	}
)

func New() *cli.Command {
	run := func(c *cli.Context) error {
		ctx := context.Background()
		scope := gdrive.DriveScope
		_, err := drive.GetClient(ctx, scope, credFlag, cacheFlag)
		return err
	}
	return &cli.Command{
		Name:   "token",
		Usage:  "get google oauth token",
		Flags:  flags,
		Action: run,
	}
}
