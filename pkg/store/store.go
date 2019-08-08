package store

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/timshannon/bolthold"
	"github.com/twistedogic/doom/pkg/helper"
	bolt "go.etcd.io/bbolt"
)

type Store struct {
	*bolthold.Store
}

func New(filename string) (*Store, error) {
	store, err := bolthold.Open(filename, os.ModePerm, nil)
	return &Store{store}, err
}

func (s *Store) UpsertItem(item interface{}) error {
	i := structs.New(item)
	m := i.Map()
	typeName := i.Name()
	for k, v := range m {
		if strings.ToLower(k) == "id" {
			val := fmt.Sprintf("%s-%v", typeName, v)
			return s.Bolt().Batch(func(tx *bolt.Tx) error {
				return s.TxUpsert(tx, val, item)
			})
		}
	}
	return fmt.Errorf("no ID field found for %#v", item)
}

func (s *Store) BulkUpsert(i interface{}) error {
	items, ok := helper.InterfaceToSlice(i)
	if !ok {
		return fmt.Errorf("%#v is not slice", i)
	}
	loop := helper.NewLoop(-1)
	for j := range items {
		item := items[j]
		loop.Add(1)
		go func() {
			loop.Done(s.UpsertItem(item))
		}()
	}
	return loop.Wait()
}
