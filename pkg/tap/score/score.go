package score

import (
	"context"
	"fmt"
	"io"

	"github.com/twistedogic/doom/pkg/client"
)

const (
	baseURL = "https://newnetty1.xscores.com"
	path    = "/stream?s=1&seq=0"
)

/*
type Result struct {
	Home int
	Away int
}

type Score struct {
	MatchID     int
	LastUpdate  time.Time
	Home        string
	Away        string
	FirstHalf   Result
	FullTime    Result
	ExtraTime   Result
	PenaltyKick Result
}
*/

type Client struct {
	Base string
	client.Client
}

func New(base string, rate int) Client {
	c := client.New(rate)
	return Client{base, c}
}

func (c Client) Update(ctx context.Context, w io.WriteCloser) error {
	defer w.Close()
	errCh := make(chan error)
	go func() {
		for {
			u := fmt.Sprintf("%s%s", c.Base, path)
			if err := c.GetResponse(u, w); err != nil {
				errCh <- err
				return
			}
		}
	}()
	go func() {
		<-ctx.Done()
		errCh <- fmt.Errorf("Cancel")
	}()
	return <-errCh
}
