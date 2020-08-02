package client

import (
	"context"
	"io"
	"log"
	"net/http"

	"go.uber.org/ratelimit"
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

func (c Client) request(ctx context.Context, method, url string, body io.Reader, w io.Writer) error {
	c.Take()
	log.Printf("%s %s", method, url)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
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

func (c Client) Get(ctx context.Context, url string, w io.Writer) error {
	return c.request(ctx, "GET", url, nil, w)
}
