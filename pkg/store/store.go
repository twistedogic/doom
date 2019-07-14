package store

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/timshannon/bolthold"
)

type Store struct {
	*bolthold.Store
}

func New(filename string) (*Store, error) {
	store, err := bolthold.Open(filename, os.ModePerm, nil)
	return &Store{store}, err
}

func (s *Store) UpsertItem(item interface{}) error {
	m := structs.New(item).Map()
	for k, v := range m {
		if strings.ToLower(k) == "id" {
			return s.Upsert(v, item)
		}
	}
	return fmt.Errorf("no ID field found for %#v", item)
}
