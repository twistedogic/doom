package store

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/structs"
	"github.com/timshannon/bolthold"
	"github.com/twistedogic/doom/pkg/helper/flatten"
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
	wg := &sync.WaitGroup{}
	items, ok := flatten.InterfaceToSlice(i)
	if !ok {
		return fmt.Errorf("%#v is not slice", i)
	}
	errCh := make(chan error)
	for j := range items {
		item := items[j]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.UpsertItem(item); err != nil {
				errCh <- err
			}
		}()
	}
	go func() {
		wg.Wait()
		errCh <- nil
	}()
	return <-errCh
}

func (s *BoltStore) GetLastUpdate() time.Time {
	var ts Timestamp
	s.Get(timestampKey, &ts)
	return time.Unix(ts.Unix, int64(0))
}
