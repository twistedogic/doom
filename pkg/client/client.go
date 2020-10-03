package client

import (
	"context"
	"io"
	"log"
	"net/http"

	"go.uber.org/ratelimit"
)

const (
	agentHeader = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"
)

func NewLimiter(rate int) ratelimit.Limiter {
	if rate > 0 {
		return ratelimit.New(rate)
	}
	return ratelimit.NewUnlimited()
}

type Client struct {
	*http.Client
	ratelimit.Limiter
}

func New(rate int) Client {
	return Client{
		&http.Client{},
		NewLimiter(rate),
	}
}

func setAgent(req *http.Request) {
	header := req.Header
	header.Set("User-Agent", agentHeader)
	req.Header = header
}

func (c Client) Request(ctx context.Context, method, url string, body io.Reader, w io.Writer) error {
	c.Take()
	log.Printf("%s %s", method, url)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	setAgent(req)
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}
	return nil
}
