package odd

import (
	"io"
	"time"
)

//TODO: nested team object
type Odd struct {
	MatchID    string
	OddID      string
	BetradarID string
	MatchTime  time.Time
	Timestamp  time.Time
	Home       string
	Away       string
	Type       string
	Outcome    string
	Value      float64
}

func Transform(r io.Reader, target io.WriteCloser) error {
	return err
}
