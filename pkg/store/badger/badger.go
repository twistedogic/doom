package badger

import (
	badger "github.com/dgraph-io/badger/v2"
)

type Store struct {
	*badger.DB
}

func New(dir string) (Store, error) {
	db, err := badger.Open(badger.DefaultOptions(dir))
	return Store{db}, err
}

func (s Store) Get(key string) ([]byte, error) {
	value := []byte(nil)
	err := s.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		if v, err := item.ValueCopy(nil); err != nil {
			return err
		} else {
			value = v
		}
		return nil
	})
	return value, err
}

func (s Store) Set(key string, value []byte) error {
	return s.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func (s Store) Scan(pattern ...string) ([]string, error) {
	keys := []string{}
	opt := badger.DefaultIteratorOptions
	if len(pattern) != 0 {
		prefix := pattern[0]
		opt.Prefix = []byte(prefix)
	}
	err := s.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(opt)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			keys = append(keys, string(key))
		}
		return nil
	})
	return keys, err
}
