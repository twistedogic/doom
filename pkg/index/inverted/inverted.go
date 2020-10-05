package inverted

import (
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/twistedogic/doom/pkg/store"
	"github.com/twistedogic/doom/proto/model"
)

func keyExists(s store.Store, key string) (bool, error) {
	keys, err := s.Scan(key)
	if err != nil {
		return false, err
	}
	for _, k := range keys {
		if k == key {
			return true, nil
		}
	}
	return false, nil
}

func valueSet(values []string) []string {
	m := make(map[string]struct{})
	for _, v := range values {
		m[v] = struct{}{}
	}
	out := make([]string, 0, len(m))
	for v := range m {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func UpdateInvertedIndex(s store.Store, key string, values ...string) error {
	ii := &model.InvertedIndex{Key: key, Values: make([]string, 0)}
	exist, err := keyExists(s, key)
	if err != nil {
		return err
	}
	if exist {
		b, err := s.Get(key)
		if err != nil {
			return err
		}
		if err := proto.Unmarshal(b, ii); err != nil {
			return err
		}
	}
	ii.Values = valueSet(append(ii.GetValues(), values...))
	b, err := ii.Marshal()
	if err != nil {
		return err
	}
	return s.Set(key, b)
}
