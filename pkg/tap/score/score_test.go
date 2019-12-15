package score

import (
	"context"
	"testing"
	"time"

	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

func TestUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	ts := testutil.Setup(t, testdataPath)
	defer ts.Close()
	target := testutil.NewMockTarget(t, false, false)
	f := New(ts.URL, -1)
	f.Update(ctx, target)
	got := string(target.Bytes())
	if len(got) == 0 {
		t.Fatal(got)
	}
	t.Log(got)
}
