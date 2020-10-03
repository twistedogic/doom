package index

import (
	"github.com/twistedogic/doom/pkg/store"
)

const (
	Prefix = "index_"
)

type Index interface {
	Reindex(string, store.Store) error
	Search(string, string) ([]string, error)
	Update([]byte) error
}