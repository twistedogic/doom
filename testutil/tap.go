package testutil

import (
	"context"
	"testing"

	"github.com/twistedogic/doom/pkg/tap"
)

type setupTapTest func(*testing.T, string) tap.Tap

func TapTest(t *testing.T, testdataPath string, setup setupTapTest, log bool) {
	store := NewMockStore(t, make(map[string][]byte), false)
	ts := Setup(t, testdataPath)
	defer ts.Close()
	tt := setup(t, ts.URL)
	if err := tt.Update(context.TODO(), store); err != nil {
		t.Fatal(err)
	}
	content := store.Content()
	if len(content) == 0 {
		t.Fatal("no entry stored")
	}
	for k, v := range store.Content() {
		if len(v) == 0 {
			t.Fatalf("key %s got empty value", k)
		}
		if log {
			t.Log(k, v)
		}
	}
}
