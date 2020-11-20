package jockey

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/twistedogic/doom/proto/model"
	"github.com/twistedogic/doom/proto/source/jockey"
)

const (
	source = "jockey"

	dateFormat      = "2006-01-02"
	dateRegexFormat = `^(\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01]))*`
)

var dateMatcher = regexp.MustCompile(dateRegexFormat)

func parseDate(date string) (time.Time, error) {
	tokens := dateMatcher.FindString(date)
	return time.Parse(dateFormat, tokens)
}

func parseOdd(odd string) (float64, error) {
	parsed := strings.Split(odd, "@")
	if len(parsed) == 2 {
		return strconv.ParseFloat(parsed[1], 64)
	}
	return strconv.ParseFloat(odd, 64)
}

func modelOddFactory(id, matchId, betType string) func(string, float64) *model.Odd {
	return func(outcome string, odd float64) *model.Odd {
		return &model.Odd{
			Id:      id,
			Source:  source,
			Type:    betType,
			MatchId: matchId,
			Odd:     float32(odd),
			Outcome: outcome,
		}
	}
}

func transformHadOdd(matchId string, odd *jockey.HadOdd) []*model.Odd {
	toModelOdd := modelOddFactory(odd.GetId(), matchId, "had")
	odds := make([]*model.Odd, 0, 3)
	if v, err := parseOdd(odd.GetHome()); err == nil {
		odds = append(odds, toModelOdd("H", v))
	}
	if v, err := parseOdd(odd.GetAway()); err == nil {
		odds = append(odds, toModelOdd("A", v))
	}
	if v, err := parseOdd(odd.GetDraw()); err == nil {
		odds = append(odds, toModelOdd("D", v))
	}
	return odds
}

func transformTotalOdd(matchId string, odd *jockey.TotalOdd) []*model.Odd {
	toModelOdd := modelOddFactory(odd.GetId(), matchId, "ttg")
	odds := make([]*model.Odd, 0, 8)
	if v, err := parseOdd(odd.GetP0()); err == nil {
		odds = append(odds, toModelOdd("P0", v))
	}
	if v, err := parseOdd(odd.GetP1()); err == nil {
		odds = append(odds, toModelOdd("P1", v))
	}
	if v, err := parseOdd(odd.GetP2()); err == nil {
		odds = append(odds, toModelOdd("P2", v))
	}
	if v, err := parseOdd(odd.GetP3()); err == nil {
		odds = append(odds, toModelOdd("P3", v))
	}
	if v, err := parseOdd(odd.GetP4()); err == nil {
		odds = append(odds, toModelOdd("P4", v))
	}
	if v, err := parseOdd(odd.GetP5()); err == nil {
		odds = append(odds, toModelOdd("P5", v))
	}
	if v, err := parseOdd(odd.GetP6()); err == nil {
		odds = append(odds, toModelOdd("P6", v))
	}
	if v, err := parseOdd(odd.GetM7()); err == nil {
		odds = append(odds, toModelOdd("M7", v))
	}
	return odds
}

func transformTeam(team *jockey.Team) *model.Team {
	return &model.Team{
		Id:   team.GetId(),
		Name: team.GetName(),
	}
}

func transformScore(s *jockey.Score) (*model.Score, error) {
	home, err := strconv.ParseInt(s.GetHome(), 10, 64)
	if err != nil {
		return nil, err
	}
	away, err := strconv.ParseInt(s.GetAway(), 10, 64)
	if err != nil {
		return nil, err
	}
	return &model.Score{
		Home: home,
		Away: away,
	}, nil
}

func transformMatch(m *jockey.Match, out *model.Match) error {
	matchTime, err := parseDate(m.GetMatchTime())
	if err != nil {
		return err
	}
	matchDate, err := types.TimestampProto(matchTime)
	if err != nil {
		return err
	}
	id := m.GetId()
	home := transformTeam(m.GetHomeTeam())
	away := transformTeam(m.GetAwayTeam())
	scores := make([]*model.Score, len(m.GetScore()))
	for i, s := range m.GetScore() {
		score, err := transformScore(s)
		if err != nil {
			return err
		}
		scores[i] = score
	}
	out.Id = id
	out.MatchDate = matchDate
	out.Home = home
	out.Away = away
	out.Score = scores
	return nil
}

func TransformJockey(b []byte, out *model.Match) error {
	m := new(jockey.Match)
	if err := json.Unmarshal(b, m); err != nil {
		return err
	}
	id := m.GetId()
	fhaodds := transformHadOdd(id, m.GetFhaodds())
	hadodds := transformHadOdd(id, m.GetHadodds())
	ttgodds := transformTotalOdd(id, m.GetTtgodds())
	odds := append(fhaodds, hadodds...)
	odds = append(odds, ttgodds...)
	if err := transformMatch(m, out); err != nil {
		return err
	}
	out.Odds = odds
	return nil
}
