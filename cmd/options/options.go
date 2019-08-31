package options

import (
	"github.com/twistedogic/doom/pkg/helper/client"
	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/pkg/tap/jc"
	"github.com/twistedogic/doom/pkg/tap/radar"
	"github.com/twistedogic/doom/pkg/target"
	"github.com/twistedogic/doom/pkg/target/csv"
	"github.com/twistedogic/doom/pkg/target/ndjson"
)

var TapOptions = make(map[string]tap.Tap)
var TargetOptions = make(map[string]target.Target)

func init() {
	TapOptions["radar"] = &radar.Client{&client.Client{}}
	TapOptions["jc"] = &jc.Client{&client.Client{}}

	TargetOptions["csv"] = &csv.Target{}
	TargetOptions["ndjson"] = &ndjson.Target{}
}
