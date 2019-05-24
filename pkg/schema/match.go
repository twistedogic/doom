package schema

import (
	"time"
)

type Match struct {
	ID             int
	HomeTeamID     int
	AwayTeamID     int
	HomeScore      int
	AwayScore      int
	TournamentID   int
	HomePossession float64
	AwayPossession float64
	MatchDate      time.Time
}
