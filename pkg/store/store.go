package store

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/timshannon/bolthold"
	"github.com/twistedogic/doom/pkg/helper"
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

func (s *Store) BulkUpsert(v interface{}) error {
	items := helper.FlattenDeep(v)
	for _, item := range items {
		if err := s.UpsertItem(item); err != nil {
			return err
		}
	}
	return nil
}
