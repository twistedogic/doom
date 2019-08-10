package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/twistedogic/jsonpath"
)

const (
	timePattern = "2006-01-02T15:04:05-07:00"
)

type Odd struct {
	ID         string
	OfficialID string `jsonpath:"$.matchIDinofficial"`
	OddID      string
	MatchID    string  `jsonpath:"$.matchID"`
	RaderID    string  `jsonpath:"$.liveEvent.matchIDbetradar" prometheus:"radar"`
	Home       string  `jsonpath:"$.homeTeam.teamNameEN" prometheus:"home"`
	Away       string  `jsonpath:"$.awayTeam.teamNameEN" prometheus:"away"`
	League     string  `jsonpath:"$.league.leagueNameEN"`
	Type       string  `prometheus:"type"`
	Outcome    string  `prometheus:"outcome"`
	MinBet     float64 `prometheus:"minbet"`
	Odd        float64 `prometheus:",value"`
	MatchTime  time.Time
	LastUpdate time.Time
}

type Odds []Odd

func (o *Odds) UnmarshalJSON(b []byte) error {
	var in interface{}
	if err := json.Unmarshal(b, &in); err != nil {
		return err
	}
	return parse(in, o)
}

func parse(i interface{}, o *Odds) error {
	odds := []Odd{}
	m := i.(map[string]interface{})
	var matchTime time.Time
	var lastUpdate time.Time
	var betType string
	for k, v := range m {
		switch {
		case k == "matchTime":
			val, err := time.Parse(timePattern, v.(string))
			if err != nil {
				return err
			}
			matchTime = val
		case k == "statuslastupdated":
			val, err := time.Parse(timePattern, v.(string))
			if err != nil {
				return err
			}
			lastUpdate = val
		case strings.HasSuffix(k, "odds"):
			betType = k
			val, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%#v is not map[string]interface{}", v)
			}
			o, err := parseOdd(val)
			if err != nil {
				return err
			}
			odds = append(odds, o...)
		}
	}
	for j, v := range odds {
		if err := jsonpath.ParseJsonpath(i, &v); err != nil {
			return err
		}
		v.MatchTime, v.LastUpdate, v.Type = matchTime, lastUpdate, betType
		v.ID = fmt.Sprintf("%s_%s_%s", v.OfficialID, v.Type, v.Outcome)
		odds[j] = v
	}
	*o = odds
	return nil
}

func parseOdd(in map[string]interface{}) ([]Odd, error) {
	var odds []Odd
	var id string
	for k, i := range in {
		v := i.(string)
		if k == "ID" {
			id = v
		}
		tokens := strings.Split(v, "@")
		if len(tokens) == 2 {
			minV, oddV := tokens[0], tokens[1]
			min, err := strconv.ParseFloat(minV, 64)
			if err != nil {
				return odds, err
			}
			odd, err := strconv.ParseFloat(oddV, 64)
			if err != nil {
				return odds, err
			}
			odds = append(odds, Odd{Outcome: k, MinBet: min, Odd: odd})
		}
	}
	for i, v := range odds {
		v.OddID = id
		odds[i] = v
	}
	return odds, nil
}
