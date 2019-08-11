package radar

import (
	"fmt"

	"github.com/twistedogic/doom/pkg/helper"
	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/pkg/target"
	"github.com/twistedogic/jsonpath"
	"go.uber.org/ratelimit"
)

const (
	maxOffset  = -180
	RadarURL   = "https://lsc.fn.sportradar.com/hkjc/en"
	MatchPath  = "$.doc[*].data[*].realcategories[*].tournaments[*].matches"
	DetailPath = "$.doc[*].data[*]"
)

type Client struct {
	ratelimit.Limiter
	BaseURL string
}

func New(u string, rate int) *Client {
	limiter := helper.NewLimiter(rate)
	return &Client{limiter, u}
}

func (c *Client) GetMatchFullFeed(offset int, value interface{}) error {
	c.Take()
	u := fmt.Sprintf("%s/Asia:Shanghai/gismo/event_fullfeed/%d/55", c.BaseURL, offset)
	return helper.GetJSON(u, value)
}

func (c *Client) GetMatchDetail(matchID int, value interface{}) error {
	c.Take()
	u := fmt.Sprintf("%s/Etc:UTC/gismo/match_details/%d", c.BaseURL, matchID)
	return helper.GetJSON(u, value)
}

func (c *Client) GetLastMatches(teamID int, value interface{}) error {
	c.Take()
	u := fmt.Sprintf("%s/Etc:UTC/gismo/stats_team_lastx/%d/5", c.BaseURL, teamID)
	return helper.GetJSON(u, value)
}

func (c *Client) GetBet(offset int, value interface{}) error {
	c.Take()
	u := fmt.Sprintf("%s/Asia:Shanghai/gismo/bet_get/hkjc/%d", c.BaseURL, offset)
	return helper.GetJSON(u, value)
}

func (c *Client) GetMatch(offset int) ([]model.Match, error) {
	var container interface{}
	if err := c.GetMatchFullFeed(offset, &container); err != nil {
		return nil, err
	}
	items, err := helper.ExtractJsonPath(container, MatchPath)
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
	items, err := helper.ExtractJsonPath(container, DetailPath)
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

func (c *Client) Update(t target.Target) error {
	for offset := 0; offset >= maxOffset; offset-- {
		items, err := c.GetMatch(offset)
		if err != nil {
			return err
		}
		for i := range items {
			id := items[i].ID
			details, err := c.GetDetail(id)
			if err != nil {
				return err
			}
			if err := t.BulkUpsert(details); err != nil {
				return err
			}
		}
	}
	return nil
}
