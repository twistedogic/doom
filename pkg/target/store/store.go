package store

import (
	"os"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/timshannon/bolthold"
	"github.com/twistedogic/doom/pkg/config"
	bolt "go.etcd.io/bbolt"
)

type BoltStore struct {
	Store *bolthold.Store
	Path  string
}

func New(filename string) (*BoltStore, error) {
	store, err := bolthold.Open(filename, os.ModePerm, nil)
	return &BoltStore{store, filename}, err
}

func (s *BoltStore) Load(c config.Setting) error {
	if err := c.ParseConfig(s); err != nil {
		return err
	}
	store, err := bolthold.Open(s.Path, os.ModePerm, nil)
	if err != nil {
		return err
	}
	s.Store = store
	return nil
}

func (s *BoltStore) UpsertItem(item interface{}) error {
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	id := uuid.NewMD5(uuid.Nil, b)
	return s.Store.Bolt().Batch(func(tx *bolt.Tx) error {
		return s.Store.TxUpsert(tx, id, item)
	})
}
