package schedule

import (
	"io"
	"log"
	"net/http"

	json "github.com/json-iterator/go"
	"github.com/robfig/cron"
)

type Scheduler struct {
	*cron.Cron
}

func New() *Scheduler {
	return &Scheduler{cron.New()}
}

func (s *Scheduler) Report(w io.Writer) error {
	entries := s.Entries()
	return json.NewEncoder(w).Encode(entries)
}

func (s *Scheduler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	if err := s.Report(res); err != nil {
		log.Println(err)
		http.Error(res, err.Error(), 500)
	}
}
