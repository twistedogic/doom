package store

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"strings"

	"github.com/twistedogic/doom/pkg/index"
	"github.com/twistedogic/doom/pkg/store"
)

const (
	sep = "|"
)

func getFieldKey(field string) string {
	return fmt.Sprintf("%s%s", index.Prefix, strings.ToLower(field))
}

func getTermEntry(term, key string) string {
	return strings.Join(
		[]string{
			strings.ToLower(term),
			strings.ToLower(key),
		},
		sep,
	)
}

func walkIndex(b []byte, idx int) string {
	cursor := idx
	s := []byte(sep)
	step := 5
	for {
		if cursor+step > len(b) {
			tokens := bytes.Split(b[idx:], s)
			return string(tokens[len(tokens)-1])
		}
		slice := b[idx : cursor+step]
		if tokens := bytes.Split(slice, s); len(tokens) > 2 {
			return string(tokens[1])
		}
		cursor += step
	}
	return ""
}

type Index struct {
	store.Store
}

func New(s store.Store) Index {
	return Index{s}
}

func (i Index) updateSuffixarray(entry FieldTermPair, key string) error {
	field := getFieldKey(strings.ToLower(entry.Field))
	term := getTermEntry(entry.Term, key)
	b, err := i.Get(field)
	if err != nil {
		return err
	}
	value := [][]byte{b, []byte(term)}
	return i.Set(field, bytes.Join(value, []byte(sep)))
}

func (i Index) Update(key string, b []byte) error {
	fields, err := ParseFieldTermPairs(b)
	if err != nil {
		return err
	}
	for _, field := range fields {
		if err := i.updateSuffixarray(field, key); err != nil {
			return err
		}
	}
	return nil
}

func (i Index) Reindex(s store.Store) error {
	keys, err := s.Scan()
	if err != nil {
		return err
	}
	for _, key := range keys {
		item, err := i.Get(key)
		if err != nil {
			return err
		}
		if err := i.Update(key, item); err != nil {
			return err
		}
	}
	return nil
}

func (i Index) Search(field, term string) ([]string, error) {
	fieldKey := getFieldKey(field)
	b, err := i.Get(fieldKey)
	if err != nil {
		return nil, err
	}
	index := suffixarray.New(b)
	indices := index.Lookup([]byte(term), -1)
	keys := make([]string, len(indices))
	for i, idx := range indices {
		keys[i] = walkIndex(b, idx)
	}
	return keys, nil
}
