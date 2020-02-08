package testutil

import (
	"bytes"
	"testing"

	"github.com/twistedogic/doom/pkg/store"
)

type entry struct {
	key   string
	value []byte
}

type StoreFactory interface {
	Setup() store.Store
	Cleanup()
}

func StoreTest(t *testing.T, factory StoreFactory) {
	cases := map[string]struct {
		input []entry
		want  map[string][]byte
	}{
		"base": {
			input: []entry{
				{"test", []byte("test")},
			},
			want: map[string][]byte{
				"test": []byte("test"),
			},
		},
		"duplicate": {
			input: []entry{
				{"test", []byte("test")},
				{"test", []byte("new")},
			},
			want: map[string][]byte{
				"test": []byte("new"),
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			s := factory.Setup()
			defer factory.Cleanup()
			t.Run("set", func(t *testing.T) {
				for _, item := range tc.input {
					if err := s.Set(item.key, item.value); err != nil {
						t.Fatal(err)
					}
				}
			})
			t.Run("get", func(t *testing.T) {
				for key, want := range tc.want {
					got, err := s.Get(key)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(got, want) {
						t.Fatalf("key: %s, got: %s, want: %s", key, string(got), string(want))
					}
				}
			})
			t.Run("scan", func(t *testing.T) {
				got, err := s.Scan()
				if err != nil {
					t.Fatal(err)
				}
				if len(got) != len(tc.want) {
					t.Fatalf("want: %d, got: %d", len(tc.want), len(got))
				}
				for _, key := range got {
					if _, ok := tc.want[key]; !ok {
						t.Fatalf("%s is not expected", key)
					}
				}
			})
		})
	}
}
