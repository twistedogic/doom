package store

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/timshannon/bolthold"
	"github.com/twistedogic/doom/pkg/helper"
	bolt "go.etcd.io/bbolt"
)

const timestampKey = "lastupdate"

type Timestamp struct {
	Unix int64
}

type BoltStore struct {
	*bolthold.Store
}

func New(filename string) (*BoltStore, error) {
	store, err := bolthold.Open(filename, os.ModePerm, nil)
	return &BoltStore{store}, err
}

func (s *BoltStore) UpdateTimestamp() error {
	now := Timestamp{time.Now().Unix()}
	return s.Bolt().Batch(func(tx *bolt.Tx) error {
		return s.TxUpsert(tx, timestampKey, now)
	})
}

func (s *BoltStore) UpsertItem(item interface{}) error {
	defer s.UpdateTimestamp()
	i := structs.New(item)
	m := i.Map()
	for k, v := range m {
		if strings.ToLower(k) == "id" {
			return s.Bolt().Batch(func(tx *bolt.Tx) error {
				return s.TxUpsert(tx, v, item)
			})
		}
	}
	return fmt.Errorf("no ID field found for %#v", item)
}

func (s *BoltStore) BulkUpsert(i interface{}) error {
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

func (s *BoltStore) GetLastUpdate() time.Time {
	var ts Timestamp
	s.Get(timestampKey, &ts)
	return time.Unix(ts.Unix, int64(0))
}
