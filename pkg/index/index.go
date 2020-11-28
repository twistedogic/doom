package index

import (
	"fmt"
	"time"

	"github.com/twistedogic/doom/pkg/store"
)

type Query struct {
	Home, Away string
	Start, End time.Time
}

func (q Query) Validate() error {
	switch {
	case q.End.IsZero() && q.Start.IsZero():
		return nil
	case q.End.IsZero():
		return fmt.Errorf("End is not set")
	case q.Start.IsZero():
		return fmt.Errorf("Start is not set")
	case q.End.Before(q.Start):
		return fmt.Errorf("End is before Start")
	}
	return nil
}

type Index interface {
	store.Setter
	Search(Query) ([]string, error)
}
