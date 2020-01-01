package score

import (
	"fmt"
	"time"
)

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

func Parse(b []byte) ([]Score, error) {
	return nil, fmt.Errorf("not implemented")
}
