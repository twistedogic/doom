package radar

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

func compare(t *testing.T, want []int, ch chan int) {
	t.Helper()
	got := []int{}
	for i := range ch {
		got = append(got, i)
	}
	sort.Ints(got)
	sort.Ints(want)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestGetMatchID(t *testing.T) {
	f, err := os.Open(filepath.Join(testdataPath, "fullfeed"))
	outCh := make(chan int)
	errCh := make(chan error)
	want := []int{
		14701401, 14701403, 14701405, 14701407, 14701409, 14701411, 14701413, 14701415,
		14701417, 14701419, 15145935, 15145939, 15145945, 15150689, 15160367, 15160375,
		16542357, 16542361, 16542365, 16542371, 16599583, 16682153, 16682155, 16682157,
		16682159, 16682161, 16682163, 16682165, 16682167, 16682169, 16682171, 17129283,
		17129285, 17129287, 17129295, 17215103, 17215105, 17215107, 17430413, 17430415,
		17430417, 17430799, 17430801, 17430803, 17430805, 17430807, 17430809, 17466911,
		17466913, 17466923, 17466925, 17466935, 17466937, 17466947, 17466949, 17527375,
		17788896, 17788908, 17873198, 18013469, 18066039, 18220065, 18254981, 18277101,
		18277121, 18284963, 18285207,
	}
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		defer close(outCh)
		errCh <- getMatchID(f, outCh)
	}()
	go func() {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}()
	compare(t, want, outCh)
}

func TestGetFeed(t *testing.T) {
	ts := testutil.Setup(t, testdataPath)
	defer ts.Close()
	f := New(ts.URL, -1)
	target := testutil.NewMockTarget(t, false, false)
	buf := new(bytes.Buffer)
	if err := f.GetFeed(0, buf, target); err != nil {
		t.Fatal(err)
	}
	if len(buf.Bytes()) == 0 {
		t.Fail()
	}
}

func TestGetDetail(t *testing.T) {
	ts := testutil.Setup(t, testdataPath)
	defer ts.Close()
	f := New(ts.URL, -1)
	target := testutil.NewMockTarget(t, false, false)
	if err := f.GetDetail(0, target); err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	ts := testutil.Setup(t, testdataPath)
	defer ts.Close()
	ctx := context.Background()
	target := testutil.NewMockTarget(t, false, false)
	f := New(ts.URL, -1)
	if err := f.Update(ctx, target); err != nil {
		t.Fatal(err)
	}
}
