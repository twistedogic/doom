package bolt

import (
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
)

const BucketName = "data"

type Store struct {
	*bolt.DB
}

func New(path string) (Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	return Store{db}, err
}

func (s Store) Get(key string) ([]byte, error) {
	value := make([]byte, 0)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return fmt.Errorf("bucket %s missing", BucketName)
		}
		v := b.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("key %s not exists", key)
		}
		value = append(value, v...)
		return nil
	})
	return value, err
}

func (s Store) Set(key string, b []byte) error {
	return db.Batch(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), b)
	})
}

func (s Store) Scan(patterns ...string) ([]string, error) {
	keys := make([]string, 0)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return fmt.Errorf("bucket %s missing", BucketName)
		}
		return b.ForEach(func(k, v []byte) error {
			key := string(k)
			if len(patterns) == 0 {
				keys = append(keys, key)
				return nil
			}
			for _, pattern := range patterns {
				if !strings.Contains(key, patterns) {
					return nil
				}
			}
			keys = append(keys, key)
			return nil
		})
	})
	return keys, err
}
