package history

import (
	"context"
	"net/url"
	"path"
	"sync"

	"github.com/twistedogic/doom/pkg/client"
	"github.com/twistedogic/doom/pkg/client/crawl"
	"github.com/twistedogic/doom/pkg/tap"
)

const (
	Base = "https://www.football-data.co.uk/data.php"
)

func IsSameDomain(urls ...string) bool {
	base, others := urls[0], urls[1:]
	baseURL, err := url.Parse(base)
	if err != nil {
		return false
	}
	for _, u := range others {
		parsedURL, err := url.Parse(u)
		if err != nil {
			return false
		}
		if baseURL.Hostname() != parsedURL.Hostname() {
			return false
		}
	}
	return true
}

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

func (c *Client) ProcessLink(link string, target tap.Target, outCh chan string) error {
	if c.isVisited(link) || !IsSameDomain(c.BaseURL, link) {
		return nil
	}
	if path.Ext(link) == ".csv" {
		return c.WriteToTarget(link, target)
	}
	return crawl.CrawlHref(link, outCh)
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
		for l := range linkCh {
			wg.Add(1)
			go func(link string) {
				defer wg.Done()
				if err := c.ProcessLink(link, target, linkCh); err != nil {
					errCh <- err
				}
			}(l)
		}
	}()
	go func() {
		wg.Wait()
		close(linkCh)
		errCh <- nil
	}()
	for {
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
		case err := <-errCh:
			return err
		}
	}
	return nil
}
