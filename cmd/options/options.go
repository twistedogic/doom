package options

import (
	"fmt"
	"time"

	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/helper/client"
	"github.com/twistedogic/doom/pkg/schedule/job"
	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/pkg/tap/backfill"
	"github.com/twistedogic/doom/pkg/tap/jc"
	"github.com/twistedogic/doom/pkg/tap/radar"
	"github.com/twistedogic/doom/pkg/target"
	"github.com/twistedogic/doom/pkg/target/csv"
	"github.com/twistedogic/doom/pkg/target/drive"
	"github.com/twistedogic/doom/pkg/target/ndjson"
)

var TapOptions = make(map[string]tap.Tap)
var TargetOptions = make(map[string]target.Target)

func init() {
	TapOptions["radar"] = &radar.Client{&client.Client{}}
	TapOptions["jc"] = &jc.Client{&client.Client{}}
	TapOptions["backfill"] = &backfill.Tap{}

	TargetOptions["csv"] = &csv.Target{}
	TargetOptions["ndjson"] = &ndjson.Target{}
	TargetOptions["drive"] = &drive.Drive{}
}

func Load(t config.Task) (*job.Job, error) {
	timeout, err := time.ParseDuration(t.Timeout)
	if err != nil {
		return nil, err
	}
	tap, ok := TapOptions[t.Tap.Name]
	if !ok {
		return nil, fmt.Errorf("no option %s for tap", t.Tap.Name)
	}
	target, ok := TargetOptions[t.Target.Name]
	if !ok {
		return nil, fmt.Errorf("no option %s for target", t.Target.Name)
	}
	if err := tap.Load(t.Tap); err != nil {
		return nil, err
	}
	if err := target.Load(t.Target); err != nil {
		return nil, err
	}
	return job.New(t.Name, tap, target, timeout), nil
}
