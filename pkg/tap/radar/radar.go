package radar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/twistedogic/doom/pkg/client"
)

const (
	maxOffset     = -180
	RadarURL      = "https://lsc.fn.sportradar.com/hkjc/en"
	FullFeedPath  = "/Asia:Shanghai/gismo/event_fullfeed/%d/55"
	DetailPath    = "/Etc:UTC/gismo/match_details/%d"
	LastMatchPath = "/Etc:UTC/gismo/stats_team_lastx/%d/5"
	BetPath       = "/Asia:Shanghai/gismo/bet_get/hkjc/%d"
)

type Client struct {
	BaseURL string
	*client.Client
}

func New(u string, rate int) *Client {
	c := client.New(rate)
	return &Client{u, c}
}

func (c *Client) GenerateURL(path string, id int) string {
	requestPath := fmt.Sprintf(path, id)
	return fmt.Sprintf("%s%s", c.BaseURL, requestPath)
}

func (c *Client) GetFeed(offset int, w io.Writer) error {
	u := c.GenerateURL(FullFeedPath, offset)
	return c.GetResponse(u, w)
}

func (c *Client) GetDetail(id int, w io.Writer) error {
	u := c.GenerateURL(DetailPath, id)
	return c.GetResponse(u, w)
}

func getMatchID(r io.Reader, ch chan int) error {
	isMatch := false
	dec := json.NewDecoder(r)
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if key, ok := t.(string); ok && (key == "_doc" || key == "_id") {
			value, _ := dec.Token()
			switch v := value.(type) {
			case string:
				if v == "match" {
					isMatch = true
				}
			case float64:
				if isMatch {
					ch <- int(v)
					isMatch = false
				}
			}
		}
	}
	return nil
}

func (c *Client) Update(ctx context.Context, w io.WriteCloser) error {
	defer w.Close()
	wg := &sync.WaitGroup{}
	matchIDCh := make(chan int)
	errCh := make(chan error)
	for i := maxOffset; i <= 0; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			buf := &bytes.Buffer{}
			if err := c.GetFeed(offset, buf); err != nil {
				errCh <- err
				return
			}
			tee := io.TeeReader(buf, w)
			if err := getMatchID(tee, matchIDCh); err != nil {
				errCh <- err
				return
			}
		}(i)
	}
	go func() {
		defer close(matchIDCh)
		wg.Wait()
	}()
	go func() {
		for id := range matchIDCh {
			if err := c.GetDetail(id, w); err != nil {
				errCh <- err
				return
			}
		}
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout")
	case err := <-errCh:
		return err
	}
	return nil
}
