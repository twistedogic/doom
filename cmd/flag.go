package cmd

import (
	"fmt"
	"time"

	"github.com/twistedogic/doom/pkg/schedule/job"
	"github.com/twistedogic/doom/pkg/tap/radar"
	"github.com/twistedogic/doom/pkg/target/csv"
	"github.com/urfave/cli"
)

const (
	dateNameFormat = "20060102"
)

var (
	tapFlag    string
	targetFlag string
	rateFlag   int

	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "tap, s",
			Value:       "radar",
			Usage:       "data source",
			Destination: &tapFlag,
		},
		cli.StringFlag{
			Name:        "target, d",
			Value:       "csv",
			Usage:       "data destination",
			Destination: &targetFlag,
		},
		cli.IntFlag{
			Name:        "rate, r",
			Value:       -1,
			Usage:       "data ingestion rate (entry per second)",
			Destination: &rateFlag,
		},
	}
)

func Run() cli.Command {
	run := func(c *cli.Context) error {
		j := job.New()
		switch tapFlag {
		case "radar":
			j.SetSrc(radar.New(radar.RadarURL, rateFlag))
		default:
			return fmt.Errorf("invalid tap option: %s", tapFlag)
		}
		switch targetFlag {
		case "csv":
			date := time.Now().Format(dateNameFormat)
			filename := fmt.Sprintf("%s_%s_%s.csv", tapFlag, targetFlag, date)
			dst, err := csv.New(filename)
			if err != nil {
				return err
			}
			j.SetDst(dst)

		default:
			return fmt.Errorf("invalid target option: %s", targetFlag)
		}
		return j.Run()
	}
	return cli.Command{
		Name:   "run",
		Usage:  "run data pipeline",
		Flags:  flags,
		Action: run,
	}
}
