package backfill

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/fatih/structs"
	"github.com/google/go-cmp/cmp"
	"github.com/twistedogic/doom/pkg/tap/backfill/model"
)

const testdataPath = "../../../testdata"

func write(t *testing.T, path string, r io.Reader) {
	t.Helper()
	filePath := filepath.Join(testdataPath, path)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filePath, b, 0644); err != nil {
		t.Fatal(err)
	}
}

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

func TestFetchPHPLinks(t *testing.T) {
	cases := map[string]struct {
		php  []byte
		want []string
	}{
		"base": {
			php: []byte(`<html>
				<a href="a.php">a</a>
				<a href="a.csv">a</a>
			</html>`),
			want: []string{"a.php"},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			ts := setup(t, "csv", tc.php)
			defer ts.Close()
			tap := New(ts.URL, -1)
			ch := make(chan string)
			go func() {
				defer close(ch)
				if err := tap.FetchPHPLinks(ts.URL, ch); err != nil {
					t.Fatal(err)
				}
			}()
			compare(t, ts.URL, ch, tc.want)
		})

	}
}

func TestGetEntry(t *testing.T) {
	php := []byte(`<html>
		<a href="a.php">a</a>
		<a href="a.csv">a</a>
	</html>`)
	ts := setup(t, "csv", php)
	defer ts.Close()
	tap := New(ts.URL, -1)
	ch := make(chan model.Entry)
	go func() {
		if err := tap.GetEntry(ch); err != nil {
			t.Fatal(err)
		}
	}()
	for entry := range ch {
		if structs.IsZero(entry) {
			t.Fatal(entry)
		}
	}
}
