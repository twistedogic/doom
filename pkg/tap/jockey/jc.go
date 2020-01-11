package jockey

import (
	"context"
	"fmt"
	"io"

	"github.com/twistedogic/doom/pkg/client"
)

const (
	Base = "https://bet.hkjc.com"
	Path = "/football/getJSON.aspx?jsontype=odds_%s.aspx"
)

var BetTypes = []string{"had", "fha", "hha", "hft"}

type Client struct {
	BaseURL string
	BetType string
	client.Client
}

func New(u, bettype string, rate int) Client {
	c := client.New(rate)
	return Client{u, bettype, c}
}

func (c Client) generateURL() string {
	requestPath := fmt.Sprintf(Path, c.BetType)
	return fmt.Sprintf("%s%s", c.BaseURL, requestPath)
}

func (c Client) Update(ctx context.Context, w io.Writer) error {
	errCh := make(chan error)
	go func() {
		errCh <- c.GetResponse(c.generateURL(), w)
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("timeout")
	}
	return nil
}
