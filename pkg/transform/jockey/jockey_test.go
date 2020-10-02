package jockey

import (
	"context"
	"testing"
	"time"

	"github.com/twistedogic/doom/pkg/tap/jockey"
	"github.com/twistedogic/doom/pkg/transform"
	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

func Test_TransformJockey(t *testing.T) {
	store := testutil.NewMockStore(t, make(map[string][]byte), false)
	ts := testutil.Setup(t, testdataPath)
	defer ts.Close()
	tap := jockey.NewResultTap(ts.URL, -1, time.Now().Add(-24*time.Hour), time.Now())
	tt := transform.New(TransformJockey, store)
	if err := tap.Update(context.TODO(), tt); err != nil {
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
		t.Log(k)
	}
}
