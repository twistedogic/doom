package backfill

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/helper/client"
	"github.com/twistedogic/doom/pkg/helper/crawl"
	"github.com/twistedogic/doom/pkg/tap/backfill/model"
	"github.com/twistedogic/doom/pkg/target"
	"go.uber.org/ratelimit"
)

const (
	Base = "https://www.football-data.co.uk/data.php"
)

type Tap struct {
	BaseURL string
	Rate    int
	visited *sync.Map
	Limiter ratelimit.Limiter
}

func New(u string, rate int) *Tap {
	return &Tap{
		BaseURL: u,
		Rate:    rate,
		visited: &sync.Map{},
		Limiter: client.NewLimiter(rate),
	}
}

func (t *Tap) CleanVisited() {
	t.visited = &sync.Map{}
}

func (t *Tap) Load(s config.Setting) error {
	if err := s.ParseConfig(t); err != nil {
		return err
	}
	t.Limiter = client.NewLimiter(t.Rate)
	t.CleanVisited()
	return nil
}

func (t *Tap) FetchCSV(u string, ch chan model.Entry) error {
	if _, loaded := t.visited.LoadOrStore(u, true); loaded {
		return nil
	}
	t.Limiter.Take()
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return model.NewDecoder(res.Body).Decode(ch)
}

func (t *Tap) FetchCSVLinks(u string, ch chan string) error {
	if _, loaded := t.visited.LoadOrStore(u, true); loaded {
		return nil
	}
	csvLinks := make(chan string)
	go func() {
		for link := range csvLinks {
			if strings.HasSuffix(link, ".csv") {
				ch <- link
			}
		}
	}()
	if err := crawl.CrawlHref(u, csvLinks); err != nil {
		return err
	}
	return nil
}

func (t *Tap) FetchPHPLinks(u string, ch chan string) error {
	phpLinks := make(chan string)
	go func() {
		for link := range phpLinks {
			if strings.HasSuffix(link, ".php") {
				ch <- link
			}
		}
	}()
	if err := crawl.CrawlHref(u, phpLinks); err != nil {
		return err
	}
	return nil
}

func (t *Tap) GetEntry(ch chan model.Entry) error {
	t.CleanVisited()
	phpLinks := make(chan string)
	csvLinks := make(chan string)
	errCh := make(chan error)
	go func() {
		defer close(phpLinks)
		if err := t.FetchPHPLinks(t.BaseURL, phpLinks); err != nil {
			log.Println(err)
			errCh <- err
		}
	}()
	go func() {
		defer close(csvLinks)
		for link := range phpLinks {
			if err := t.FetchCSVLinks(link, csvLinks); err != nil {
				errCh <- err
			}
		}
	}()
	go func() {
		for link := range csvLinks {
			if err := t.FetchCSV(link, ch); err != nil {
				errCh <- err
			}
		}
		errCh <- nil
	}()
	return <-errCh
}

func (t *Tap) Update(dst target.Target) error {
	ch := make(chan model.Entry)
	errCh := make(chan error)
	go func() {
		errCh <- t.GetEntry(ch)
	}()
	go func() {
		for entry := range ch {
			if err := dst.UpsertItem(entry); err != nil {
				errCh <- err
			}
		}
	}()
	return <-errCh
}
