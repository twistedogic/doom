package store

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/timshannon/bolthold"
)

func SetFile(t *testing.T) (string, func()) {
	t.Helper()
	f, err := ioutil.TempFile("", "store")
	if err != nil {
		t.Fatal(err)
	}
	return f.Name(), func() {
		f.Close()
		os.Remove(f.Name())
	}
}

type Nested struct {
	Name  string
	Value int
}

type Item struct {
	ID     uint64
	Name   string
	Value  int
	Nested Nested
}

func TestStore(t *testing.T) {
	cases := map[string]struct {
		data  []Item
		field string
		query string
		want  []Item
	}{
		"base": {
			data: []Item{
				{ID: uint64(0), Name: "something", Value: 2, Nested: Nested{Name: "test", Value: 1}},
				{ID: uint64(1), Name: "other", Value: 2, Nested: Nested{Name: "test", Value: 1}},
			},
			field: "Name",
			query: "other",
			want: []Item{
				{ID: uint64(1), Name: "other", Value: 2, Nested: Nested{Name: "test", Value: 1}},
			},
		},
		"nested": {
			data: []Item{
				{ID: uint64(0), Name: "something", Value: 2, Nested: Nested{Name: "test", Value: 1}},
				{ID: uint64(1), Name: "other", Value: 2, Nested: Nested{Name: "test", Value: 1}},
			},
			field: "Nested.Name",
			query: "test",
			want: []Item{
				{ID: uint64(0), Name: "something", Value: 2, Nested: Nested{Name: "test", Value: 1}},
				{ID: uint64(1), Name: "other", Value: 2, Nested: Nested{Name: "test", Value: 1}},
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			filename, cleanup := SetFile(t)
			defer cleanup()
			s, err := New(filename)
			if err != nil {
				t.Fatal(err)
			}
			for _, item := range tc.data {
				if err := s.UpsertItem(item); err != nil {
					t.Fatal(err)
				}
			}
			got := []Item{}
			if err := s.Find(&got, bolthold.Where(tc.field).Eq(tc.query)); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Fatalf("want %#v, got %#v", tc.want, got)
			}
		})
	}
}
