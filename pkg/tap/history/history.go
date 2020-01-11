package history

import (
	"context"
	"io"
	"net/http"
	"path"
	"sync"

	"github.com/twistedogic/doom/pkg/client/crawl"
)

const (
	Base = "https://www.football-data.co.uk/data.php"
)

type Client struct {
	sync.Mutex
	BaseURL string
	visited *sync.Map
}

func New(u string) *Client {
	return &Client{BaseURL: u}
}

func (c *Client) isVisited(u string) bool {
	if c.visited == nil {
		c.visited = &sync.Map{}
	}
	_, loaded := c.visited.LoadOrStore(u, true)
	return loaded
}

func (c *Client) clearVisited() {
	c.visited = &sync.Map{}
}

func (c *Client) GetCSV(u string, w io.Writer) error {
	c.Lock()
	defer c.Unlock()
	if c.isVisited(u) {
		return nil
	}
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}
	return nil
}

func (c *Client) FetchLink(link string, outCh chan string) error {
	errCh := make(chan error)
	linkCh := make(chan string)
	go func() {
		for l := range linkCh {
			outCh <- l
		}
		errCh <- nil
	}()
	go func() {
		defer close(linkCh)
		if !c.isVisited(link) {
			if err := crawl.CrawlHref(link, linkCh); err != nil {
				errCh <- err
			}
		}
	}()
	return <-errCh
}

func (c *Client) Update(ctx context.Context, w io.Writer) error {
	c.clearVisited()
	wg := new(sync.WaitGroup)
	errCh := make(chan error)
	linkCh := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()
		linkCh <- c.BaseURL
	}()
	go func() {
		<-ctx.Done()
		errCh <- ctx.Err()
	}()
	go func() {
		wg.Wait()
		close(linkCh)
		errCh <- nil
	}()
	for l := range linkCh {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			switch path.Ext(link) {
			case ".csv":
				if err := c.GetCSV(link, w); err != nil {
					errCh <- err
				}
			default:
				if err := c.FetchLink(link, linkCh); err != nil {
					errCh <- err
				}
			}
		}(l)
	}
	return <-errCh
}
