package radar

import (
	"fmt"

	"github.com/twistedogic/doom/pkg/client"
)

const (
	RadarURL = "https://lsc.fn.sportradar.com/hkjc/en"
)

type Client struct {
	BaseURL string
}

func New(u string) *Client {
	return &Client{u}
}

func (c *Client) GetMatchFullFeed(offset int, value interface{}) error {
	u := fmt.Sprintf("%s/Asia:Shanghai/gismo/event_fullfeed/%d/55", c.BaseURL, offset)
	return client.GetJSON(u, value)
}

func (c *Client) GetMatchDetail(matchID int, value interface{}) error {
	u := fmt.Sprintf("%s/Etc:UTC/gismo/match_details/%d", c.BaseURL, matchID)
	return client.GetJSON(u, value)
}

func (c *Client) GetLastMatches(teamID int, value interface{}) error {
	u := fmt.Sprintf("%s/Etc:UTC/gismo/stats_team_lastx/%d/5", c.BaseURL, matchID)
	return client.GetJSON(u, value)
}

//https://lsc.fn.sportradar.com/hkjc/en/Etc:UTC/gismo/stats_team_lastx/39/5
//https://lsc.fn.sportradar.com/hkjc/en/Etc:UTC/gismo/stats_season_uniqueteamstats/54785
//https://lsc.fn.sportradar.com/hkjc/en/Etc:UTC/gismo/stats_team_versusrecent/14/34
