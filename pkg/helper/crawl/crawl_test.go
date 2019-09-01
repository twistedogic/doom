package crawl

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Setup(t *testing.T, payload []byte) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write(payload)
	}))
}

func TestCrawl(t *testing.T) {
	cases := map[string]struct {
		payload []byte
		want    []string
	}{
		"base": {
			[]byte(`<html>
				<a href="https://localhost">ok</a>
				<a href="http://localhost">ok</a>
			</html>`),
			[]string{"https://localhost", "http://localhost"},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			ts := Setup(t, tc.payload)
			defer ts.Close()
			ch := make(chan string)
			go func() {
				if err := CrawlHref(ts.URL, ch); err != nil {
					t.Fatal(err)
				}
			}()
			got := []string{}
			for u := range ch {
				got = append(got, u)
			}
			sort.Strings(got)
			sort.Strings(tc.want)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
