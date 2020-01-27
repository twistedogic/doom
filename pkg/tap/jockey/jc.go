package jockey

import (
	"context"
	"fmt"

	"github.com/twistedogic/doom/pkg/client"
	"github.com/twistedogic/doom/pkg/tap"
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

func (c Client) Update(ctx context.Context, target tap.Target) error {
	errCh := make(chan error)
	go func() {
		<-ctx.Done()
		errCh <- ctx.Err()
	}()
	go func() {
		errCh <- c.WriteToTarget(c.generateURL(), target)
	}()
	return <-errCh
}
