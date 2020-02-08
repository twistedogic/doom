package start

import (
	"time"

	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/pkg/tap/history"
	"github.com/twistedogic/doom/pkg/tap/jockey"
	"github.com/twistedogic/doom/pkg/tap/radar"
)

type Config struct {
	Name   string
	Tap    tap.Tap
	Period time.Duration
}

var configs = []Config{
	{"jc", jockey.New(jockey.Base, "had", 5), 1 * time.Minute},
	{"history", history.New(history.Base, 5), 24 * time.Hour},
	{"radar", radar.New(radar.Base, 5), 24 * time.Hour},
}
