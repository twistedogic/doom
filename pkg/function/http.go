package function

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/pkg/tap/jc"
	"github.com/twistedogic/doom/pkg/tap/jc/model"
	"github.com/twistedogic/doom/pkg/target"
	"github.com/twistedogic/doom/pkg/target/prom"
)

var (
	DefaultURL  = jc.JcURL
	DefaultRate = -1
)

type OddHTTP struct {
	target  target.Target
	tap     tap.Tap
	Handler http.Handler
}

func New() (*OddHTTP, error) {
	reg := prometheus.NewRegistry()
	target, err := prom.New(model.Odd{}, reg)
	if err != nil {
		return nil, err
	}
	tap := jc.New(DefaultURL, DefaultRate)
	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	return &OddHTTP{target, tap, handler}, nil
}

func (o *OddHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := o.tap.Update(o.target); err != nil {
		log.Print(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	o.Handler.ServeHTTP(w, r)
}
