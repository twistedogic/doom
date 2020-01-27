package history

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

func compare(t *testing.T, base string, ch chan string, expect []string) {
	t.Helper()
	got := []string{}
	want := []string{}
	for _, v := range expect {
		u := v
		if !strings.HasPrefix(v, "http") {
			u = fmt.Sprintf("%s/%s", base, strings.TrimPrefix(v, "/"))
		}
		want = append(want, u)
	}
	for v := range ch {
		got = append(got, v)
	}
	sort.Strings(got)
	sort.Strings(want)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func setup(t *testing.T, csv string, php []byte) *httptest.Server {
	t.Helper()
	filePath := filepath.Join(testdataPath, csv)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		body := php
		if strings.HasSuffix(req.URL.Path, ".csv") {
			body = b
		}
		if _, err := res.Write(body); err != nil {
			t.Fatal(err)
		}
	})
	return httptest.NewServer(handler)
}

func TestFetchLink(t *testing.T) {
	cases := map[string]struct {
		php  []byte
		want []string
	}{
		"base": {
			php: []byte(`<html>
				<a href="a.php">a</a>
				<a href="a.csv">a</a>
			</html>`),
			want: []string{"a.php", "a.csv"},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			ts := setup(t, "csv", tc.php)
			defer ts.Close()
			client := New(ts.URL, -1)
			ch := make(chan string)
			go func() {
				defer close(ch)
				if err := client.FetchLink(ts.URL, ch); err != nil {
					t.Fatal(err)
				}
			}()
			compare(t, ts.URL, ch, tc.want)
		})
	}
}

func TestUpdate(t *testing.T) {
	php := []byte(`<html>
		<a href="a.php">a</a>
		<a href="a.csv">a</a>
		<a href="b.csv">a</a>
		<a href="c.csv">a</a>
	</html>`)
	ts := setup(t, "csv", php)
	target := testutil.NewMockTarget(t, true, false)
	defer ts.Close()
	client := New(ts.URL, -1)
	ctx := context.Background()
	if err := client.Update(ctx, target); err != nil {
		t.Fatal(err)
	}
}
