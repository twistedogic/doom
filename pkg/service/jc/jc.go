package radar

import (
	"fmt"

	"github.com/twistedogic/doom/pkg/client"
)

const (
	JcURL = "https://bet.hkjc.com/football/getJSON.aspx?jsontype=odds_%s.aspx"
)

type Client struct {
	BaseURL string
}

func New(u string) *Client {
	return &Client{u}
}

func (c *Client) GetOdd(oddType string, value interface{}) error {
	u := fmt.Sprintf(c.BaseURL, oddType)
	return client.GetJSON(u, value)
}
