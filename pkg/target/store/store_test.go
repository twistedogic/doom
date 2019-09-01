package store

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/timshannon/bolthold"
	"github.com/twistedogic/doom/pkg/config"
)

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

type Items []Item

func (i Items) Len() int      { return len(i) }
func (i Items) Swap(a, b int) { i[a], i[b] = i[b], i[a] }

type byID struct{ Items }

func (b byID) Less(i, j int) bool { return b.Items[i].ID < b.Items[j].ID }

func setFile(t *testing.T) (string, func()) {
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
			filename, cleanup := setFile(t)
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
			if err := s.Store.Find(&got, bolthold.Where(tc.field).Eq(tc.query)); err != nil {
				t.Fatal(err)
			}
			sort.Sort(byID{tc.want})
			sort.Sort(byID{got})
			if !cmp.Equal(tc.want, got) {
				t.Fatalf("want %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	filename, cleanup := setFile(t)
	defer cleanup()
	cfg := config.Setting{
		Config: map[string]string{
			"path": filename,
		},
	}
	empty := &BoltStore{Store: &bolthold.Store{}}
	if err := empty.Load(cfg); err != nil {
		t.Fatal(err)
	}
}
