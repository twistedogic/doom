package function

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/pkg/tap/jc"
	"github.com/twistedogic/doom/pkg/target/prom"
)

var (
	DefaultURL  = jc.JcURL
	DefaultRate = -1
)

func OddHTTP(w http.ResponseWriter, r *http.Request) {
	tap := jc.New(DefaultURL, DefaultRate)
	target, err := prom.New(model.Odd{})
	if err != nil {
		log.Print(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	if err := tap.Update(target); err != nil {
		log.Print(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
	promhttp.Handler().ServeHTTP(w, r)
}
