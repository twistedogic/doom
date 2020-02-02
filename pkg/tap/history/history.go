package history

import (
	"context"
	"path"
	"sync"

	"github.com/twistedogic/doom/pkg/client"
	"github.com/twistedogic/doom/pkg/client/crawl"
	"github.com/twistedogic/doom/pkg/tap"
)

const (
	Base = "https://www.football-data.co.uk/data.php"
)

type Client struct {
	client.Client
	BaseURL string
	visited *sync.Map
}

func New(u string, rate int) *Client {
	c := client.New(rate)
	visited := new(sync.Map)
	return &Client{c, u, visited}
}

func (c *Client) resetVisited() {
	c.visited = new(sync.Map)
}

func (c *Client) isVisited(u string) bool {
	if c.visited == nil {
		c.resetVisited()
	}
	_, loaded := c.visited.LoadOrStore(u, true)
	return loaded
}

func (c *Client) GetCSV(u string, target tap.Target) error {
	if c.isVisited(u) {
		return nil
	}
	return c.WriteToTarget(u, target)
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

func (c *Client) Update(ctx context.Context, target tap.Target) error {
	c.resetVisited()
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
				if err := c.GetCSV(link, target); err != nil {
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
