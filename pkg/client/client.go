package client

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/twistedogic/doom/pkg/tap"
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
		return err
	}
	defer res.Body.Close()
	if _, err := io.Copy(w, res.Body); err != nil {
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

func (c Client) WriteToTarget(u string, target tap.Target) error {
	c.Take()
	log.Printf("GET %s", u)
	res, err := c.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return target.Write(b)
}
