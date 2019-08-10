package prom

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/twistedogic/doom/pkg/helper"
)

type Target struct {
	metric     *prometheus.GaugeVec
	lastupdate time.Time
}

func New(i interface{}) (*Target, error) {
	metric, err := SetMetric(i)
	if err != nil {
		return nil, err
	}
	return &Target{metric: metric}, nil
}

func (t *Target) UpsertItem(i interface{}) error {
	t.lastupdate = time.Now()
	return Update(t.metric, i)
}

func (t *Target) BulkUpsert(i interface{}) error {
	items, ok := helper.InterfaceToSlice(i)
	if !ok {
		return t.UpsertItem(i)
	}
	for _, item := range items {
		if err := t.UpsertItem(item); err != nil {
			return err
		}
	}
	return nil
}

func (t *Target) GetLastUpdate() time.Time {
	return t.lastupdate
}
