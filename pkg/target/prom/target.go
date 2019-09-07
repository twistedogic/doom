package prom

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/twistedogic/doom/pkg/config"
)

type Target struct {
	metric     *prometheus.GaugeVec
	lastupdate time.Time
}

func New(i interface{}, reg *prometheus.Registry) (*Target, error) {
	metric, err := SetMetric(i, reg)
	if err != nil {
		return nil, err
	}
	return &Target{metric: metric}, nil
}

func (t *Target) Load(s config.Setting) error {
	return fmt.Errorf("prom is long running service")
}

func (t *Target) UpsertItem(i interface{}) error {
	t.lastupdate = time.Now()
	return Update(t.metric, i)
}

func (t *Target) Close() error {
	return nil
}
