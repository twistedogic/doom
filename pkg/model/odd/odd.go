package odd

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
)

const (
	dateFormat = "2006-01-02T15:04:05-07:00"
)

type (
	MatchModel struct {
		ID         string
		BetradarID string
		MatchTime  time.Time
		Timestamp  time.Time
		Home       string
		Away       string
	}
	OddModel struct {
		ID        string
		Timestamp time.Time
		MatchID   string
		Type      string
		Outcome   string
		Value     float64
	}
)

func ParsePayout(v string) (float64, error) {
	tokens := strings.Split(v, "@")
	if len(tokens) != 2 {
		return 0, fmt.Errorf("failed to parse payout: %s", v)
	}
	return strconv.ParseFloat(tokens[1], 64)
}

func ParseOddModel(match Match) ([]OddModel, error) {
	out := []OddModel{}
	for name, field := range structs.Map(match) {
		if strings.Contains(name, "odds") {
			m, ok := field.(map[string]interface{})
			if !ok {
				continue
			}
			id, ok := m["ID"]
			if !ok {
				return nil, fmt.Errorf("no ID in %#v", m)
			}
			for k, v := range m {
				if payout, err := ParsePayout(v.(string)); err == nil {
					out = append(out, OddModel{
						ID:        id.(string),
						Timestamp: time.Now(),
						MatchID:   match.MatchID,
						Type:      name,
						Outcome:   k,
						Value:     payout,
					})
				}
			}
		}
	}
	return out, nil
}

func ParseMatchModel(m Match) (MatchModel, error) {
	match := MatchModel{
		ID:   m.MatchID,
		Home: m.HomeTeam.TeamNameEN,
		Away: m.AwayTeam.TeamNameEN,
	}
	matchTime, err := time.Parse(dateFormat, m.MatchTime)
	if err != nil {
		return match, err
	}
	match.MatchTime = matchTime
	if m.Statuslastupdated != nil {
		statusUpdate, err := time.Parse(dateFormat, *m.Statuslastupdated)
		if err != nil {
			return match, err
		}
		match.Timestamp = statusUpdate
	}
	if m.LiveEvent != nil {
		match.BetradarID = m.LiveEvent.MatchIDbetradar
	}
	return match, nil
}

func Transform(r io.Reader, target io.WriteCloser) error {
	defer target.Close()
	decoder := json.NewDecoder(r)
	encoder := json.NewEncoder(target)
	for decoder.More() {
		var odd Odd
		if err := decoder.Decode(&odd); err != nil {
			return err
		}
		for _, e := range odd {
			for _, m := range e.Matches {
				match, err := ParseMatchModel(m)
				if err != nil {
					return err
				}
				if err := encoder.Encode(match); err != nil {
					return err
				}
				odds, err := ParseOddModel(m)
				if err != nil {
					return err
				}
				for _, o := range odds {
					if err := encoder.Encode(o); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
