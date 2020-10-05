package inverted

import (
	"sort"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/twistedogic/doom/proto/model"
	"github.com/twistedogic/doom/testutil"
)

func setup(t *testing.T, indices []*model.InvertedIndex) *testutil.MockStore {
	t.Helper()
	s := testutil.NewMockStore(t, make(map[string][]byte), false)
	for _, idx := range indices {
		key := idx.GetKey()
		b, err := idx.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if err := s.Set(key, b); err != nil {
			t.Fatal(err)
		}
	}
	return s
}

func check(t *testing.T, s *testutil.MockStore, want []*model.InvertedIndex) {
	t.Helper()
	content := s.Content()
	got := make([]*model.InvertedIndex, 0, len(content))
	for _, v := range content {
		idx := new(model.InvertedIndex)
		if err := proto.Unmarshal(v, idx); err != nil {
			t.Fatal(err)
		}
		got = append(got, idx)
	}
	sort.Slice(got, func(i, j int) bool { return got[i].Key > got[j].Key })
	sort.Slice(want, func(i, j int) bool { return want[i].Key > want[j].Key })
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatal(diff)
	}
}

func Test_UpdateInvertedIndex(t *testing.T) {
	cases := map[string]struct {
		key           string
		values        []string
		indices, want []*model.InvertedIndex
	}{
		"brand new index": {
			key:    "key",
			values: []string{"value", "value1"},
			want: []*model.InvertedIndex{
				{Key: "key", Values: []string{"value", "value1"}},
			},
		},
		"new index": {
			key:    "key",
			values: []string{"value", "value1"},
			indices: []*model.InvertedIndex{
				{Key: "key1", Values: []string{"value", "value1"}},
			},
			want: []*model.InvertedIndex{
				{Key: "key", Values: []string{"value", "value1"}},
				{Key: "key1", Values: []string{"value", "value1"}},
			},
		},
		"update index": {
			key:    "key",
			values: []string{"value", "something"},
			indices: []*model.InvertedIndex{
				{Key: "key", Values: []string{"value", "value1"}},
			},
			want: []*model.InvertedIndex{
				{Key: "key", Values: []string{"something", "value", "value1"}},
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			s := setup(t, tc.indices)
			if err := UpdateInvertedIndex(s, tc.key, tc.values...); err != nil {
				t.Fatal(err)
			}
			check(t, s, tc.want)
		})
	}
}
