package start

import (
	"fmt"
	"log"
	"net/http"

	"github.com/twistedogic/doom/pkg/function"
	"github.com/urfave/cli"
)

const (
	dateNameFormat = "20060102"
)

var (
	portFlag int
	rateFlag int

	flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			Value:       3000,
			Usage:       "port for server",
			Destination: &portFlag,
		},
		cli.IntFlag{
			Name:        "rate, r",
			Value:       -1,
			Usage:       "data ingestion rate (entry per second)",
			Destination: &rateFlag,
		},
	}
)

func New() cli.Command {
	run := func(c *cli.Context) error {
		handler, err := function.New()
		if err != nil {
			return err
		}
		function.DefaultRate = rateFlag
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", portFlag),
			Handler: handler,
		}
		log.Printf("Server Running at %d", portFlag)
		return server.ListenAndServe()
	}
	return cli.Command{
		Name:   "start",
		Usage:  "start metric server",
		Flags:  flags,
		Action: run,
	}
}
