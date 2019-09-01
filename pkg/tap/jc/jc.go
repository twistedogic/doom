package jc

import (
	"fmt"
	"sync"

	"github.com/twistedogic/doom/pkg/helper/client"
	"github.com/twistedogic/doom/pkg/tap/jc/model"
	"github.com/twistedogic/doom/pkg/target"
	"github.com/twistedogic/jsonpath"
)

const (
	JcURL   = "https://bet.hkjc.com/football/getJSON.aspx?jsontype=odds_%s.aspx"
	OddPath = "$.[*].matches.[*]"
)

var BetTypes = []string{"had", "fha", "hha", "hft"}

type Client struct {
	*client.Client
}

func New(u string, rate int) *Client {
	c := client.New(u, rate)
	return &Client{c}
}

func (c *Client) GetOdd(oddType string, value interface{}) error {
	u := fmt.Sprintf(c.BaseURL, oddType)
	return c.GetJSON(u, value)
}

func (c *Client) getOdd(betType string) ([]model.Odd, error) {
	var container interface{}
	if err := c.GetOdd(betType, &container); err != nil {
		return nil, err
	}
	items, err := client.ExtractJsonPath(container, OddPath)
	if err != nil {
		return nil, err
	}
	out := []model.Odd{}
	for _, item := range items {
		var o model.Odds
		if err := jsonpath.Unmarshal(item, &o); err != nil {
			return nil, err
		}
		out = append(out, o...)
	}
	return out, nil
}

func (c *Client) GetOdds() ([]model.Odd, error) {
	wg := &sync.WaitGroup{}
	outCh := make(chan model.Odd)
	errCh := make(chan error)
	out := []model.Odd{}
	for i := range BetTypes {
		wg.Add(1)
		betType := BetTypes[i]
		go func() {
			defer wg.Done()
			odds, err := c.getOdd(betType)
			if err != nil {
				errCh <- err
				return
			}
			for _, odd := range odds {
				outCh <- odd
			}
		}()
	}
	go func() {
		defer close(outCh)
		wg.Wait()
		errCh <- nil
	}()
	for {
		select {
		case o := <-outCh:
			out = append(out, o)
		case err := <-errCh:
			return out, err
		}
	}
}

func (c *Client) Update(t target.Target) error {
	items, err := c.GetOdds()
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := t.UpsertItem(item); err != nil {
			return err
		}
	}
	return nil
}
