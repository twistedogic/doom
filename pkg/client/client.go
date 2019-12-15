package client

import (
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

func (c Client) Request(req *http.Request, w io.Writer) error {
	c.Take()
	log.Printf("%s %s", req.Method, req.URL)
	res, err := c.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	defer res.Body.Close()
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (c Client) GetResponse(u string, w io.Writer) error {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	return c.Request(req, w)
}
