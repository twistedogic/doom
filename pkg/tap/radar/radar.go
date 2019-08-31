package radar

import (
	"fmt"
	"log"
	"sync"

	"github.com/twistedogic/doom/pkg/helper/client"
	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/pkg/target"
	"github.com/twistedogic/jsonpath"
)

const (
	maxOffset  = -180
	RadarURL   = "https://lsc.fn.sportradar.com/hkjc/en"
	MatchPath  = "$.doc[*].data[*].realcategories[*].tournaments[*].matches"
	DetailPath = "$.doc[*].data[*]"
)

type Client struct {
	*client.Client
}

func New(u string, rate int) *Client {
	c := client.New(u, rate)
	return &Client{c}
}

func (c *Client) GetMatchFullFeed(offset int, value interface{}) error {
	u := fmt.Sprintf("%s/Asia:Shanghai/gismo/event_fullfeed/%d/55", c.BaseURL, offset)
	return c.GetJSON(u, value)
}

func (c *Client) GetMatchDetail(matchID int, value interface{}) error {
	u := fmt.Sprintf("%s/Etc:UTC/gismo/match_details/%d", c.BaseURL, matchID)
	return c.GetJSON(u, value)
}

func (c *Client) GetLastMatches(teamID int, value interface{}) error {
	u := fmt.Sprintf("%s/Etc:UTC/gismo/stats_team_lastx/%d/5", c.BaseURL, teamID)
	return c.GetJSON(u, value)
}

func (c *Client) GetBet(offset int, value interface{}) error {
	u := fmt.Sprintf("%s/Asia:Shanghai/gismo/bet_get/hkjc/%d", c.BaseURL, offset)
	return c.GetJSON(u, value)
}

func (c *Client) GetMatch(offset int) ([]model.Match, error) {
	var container interface{}
	if err := c.GetMatchFullFeed(offset, &container); err != nil {
		return nil, err
	}
	items, err := client.ExtractJsonPath(container, MatchPath)
	if err != nil {
		return nil, err
	}
	out := make([]model.Match, len(items))
	for i, item := range items {
		if err := jsonpath.Unmarshal(item, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (c *Client) GetDetail(matchID int) ([]model.Detail, error) {
	var container interface{}
	if err := c.GetMatchDetail(matchID, &container); err != nil {
		return nil, err
	}
	items, err := client.ExtractJsonPath(container, DetailPath)
	if err != nil {
		return nil, err
	}
	out := make([]model.Detail, len(items))
	for i, item := range items {
		if err := jsonpath.Unmarshal(item, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (c *Client) FetchMatch(offset int, matchCh chan model.Match) error {
	items, err := c.GetMatch(offset)
	if err != nil {
		return err
	}
	log.Printf("%d matches for %d", len(items), offset)
	for _, item := range items {
		if item.IsFinish() {
			matchCh <- item
		}
	}
	return nil
}

func (c *Client) FetchDetail(matchCh chan model.Match, detailCh chan model.MatchDetail) error {
	for match := range matchCh {
		details, err := c.GetDetail(match.ID)
		if err != nil {
			return err
		}
		for _, detail := range details {
			detailCh <- model.MatchDetail{match, detail}
		}
	}
	return nil
}

func (c *Client) Update(t target.Target) error {
	wg := &sync.WaitGroup{}
	matchCh := make(chan model.Match)
	detailCh := make(chan model.MatchDetail)
	errCh := make(chan error)
	for offset := 0; offset >= maxOffset; offset-- {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := c.FetchMatch(i, matchCh); err != nil {
				errCh <- err
			}
		}(offset)
	}
	go func() {
		wg.Wait()
		close(matchCh)
	}()
	go func() {
		errCh <- c.FetchDetail(matchCh, detailCh)
	}()
	var written int
	for {
		select {
		case d := <-detailCh:
			if err := t.UpsertItem(d); err != nil {
				errCh <- err
			}
			written += 1
			if written%100 == 0 {
				log.Printf("%d records written", written)
			}
		case err := <-errCh:
			return err
		}
	}
	return nil
}
