package store

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twistedogic/doom/testutil"
)

func Test_walkIndex(t *testing.T) {
	cases := map[string]struct {
		input []byte
		idx   int
		want  string
	}{
		"base": {
			input: []byte("term|key"),
			idx:   0,
			want:  "key",
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got := walkIndex(tc.input, tc.idx)
			if got != tc.want {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}

func Test_ParseFieldTermPairs(t *testing.T) {
	cases := map[string]struct {
		input []byte
		want  []FieldTermPair
	}{
		"base": {
			input: []byte(`{
"key":"value",
"slice": ["string","raw"],
"nested": {
"k": "v",
"num": 0
},
"nestedSlice": [{"message":"OK"}]
}`),
			want: []FieldTermPair{
				{"key", "value"},
				{"slice", "string"},
				{"slice", "raw"},
				{"message", "ok"},
				{"k", "v"},
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got, err := ParseFieldTermPairs(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			sort.Slice(got, func(i, j int) bool { return got[i].Field < got[j].Field })
			sort.Slice(tc.want, func(i, j int) bool { return tc.want[i].Field < tc.want[j].Field })
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

type searchCase struct {
	field string
	term  string
	want  []string
}

func Test_Index(t *testing.T) {
	cases := map[string]struct {
		content map[string][]byte
		search  []searchCase
	}{
		"base": {
			content: map[string][]byte{
				"key1": []byte(`{"key":"abc"}`),
				"key2": []byte(`{"key":"abcd"}`),
				"key3": []byte(`{"key":"abcd"}`),
			},
			search: []searchCase{
				{field: "key", term: "d", want: []string{"key2", "key3"}},
				{field: "key", term: "abc", want: []string{"key1", "key2", "key3"}},
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			store := testutil.NewMockStore(t, tc.content, false)
			index := New(store)
			if err := index.Reindex(store); err != nil {
				t.Fatal(err)
			}
			for _, search := range tc.search {
				got, err := index.Search(search.field, search.term)
				if err != nil {
					t.Fatal(err)
				}
				sort.Strings(got)
				sort.Strings(search.want)
				if diff := cmp.Diff(search.want, got); diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}
