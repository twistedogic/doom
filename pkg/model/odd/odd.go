package odd

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	json "github.com/json-iterator/go"
	"github.com/twistedogic/doom/pkg/model"
)

const (
	dateFormat = "2006-01-02T15:04:05-07:00"

	Type model.Type = "odd"
)

type (
	MatchModel struct {
		MatchID    string
		BetradarID int
		MatchTime  time.Time
		Timestamp  time.Time
		Home       string
		Away       string
	}
	Model struct {
		MatchModel
		ID      string
		Type    string
		Outcome string
		Value   float64
	}
)

func (m Model) Item(i *model.Item) error {
	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s_%s", m.Type, m.ID)
	i.Key, i.Type, i.Data = key, Type, b
	return nil
}

func parsePayout(v string) (float64, error) {
	tokens := strings.Split(v, "@")
	if len(tokens) != 2 {
		return 0, fmt.Errorf("failed to parse payout: %s", v)
	}
	return strconv.ParseFloat(tokens[1], 64)
}

func parseMatchModel(m Match) (MatchModel, error) {
	match := MatchModel{
		MatchID: m.MatchID,
		Home:    m.HomeTeam.TeamNameEN,
		Away:    m.AwayTeam.TeamNameEN,
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
	if m.LiveEvent != nil && m.LiveEvent.MatchIDbetradar != "" {
		id, err := strconv.Atoi(m.LiveEvent.MatchIDbetradar)
		if err != nil {
			return match, err
		}
		match.BetradarID = id
	}
	return match, nil
}

func parseOddModel(match Match) ([]Model, error) {
	out := []Model{}
	matchModel, err := parseMatchModel(match)
	if err != nil {
		return nil, err
	}
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
				if payout, err := parsePayout(v.(string)); err == nil {
					out = append(out, Model{
						MatchModel: matchModel,
						ID:         id.(string),
						Type:       name,
						Outcome:    k,
						Value:      payout,
					})
				}
			}
		}
	}
	return out, nil
}

func Transform(r io.Reader, encoder model.Encoder) error {
	decoder := json.NewDecoder(r)
	for decoder.More() {
		var odd Odd
		if err := decoder.Decode(&odd); err != nil {
			return err
		}
		for _, e := range odd {
			for _, m := range e.Matches {
				odds, err := parseOddModel(m)
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
