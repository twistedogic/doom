package jockey

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/twistedogic/doom/pkg/client"
	"github.com/twistedogic/doom/testutil"
)

const testdataPath = "../../../testdata"

func Test_toQueryString(t *testing.T) {
	cases := map[string]struct {
		kv   map[string]string
		want string
	}{
		"base": {
			kv: map[string]string{
				"key":   "value",
				"other": "thing",
			},
			want: "key=value&other=thing",
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got := toQueryString(tc.kv)
			if tc.want != got {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}

func Test_Store(t *testing.T) {
	cases := map[string]struct {
		path, method, typePrefix string
	}{
		"base": {
			path:       "/football/getJSON.aspx?jsontype=search_result.aspx&startdate=20180107&enddate=20180203&teamid=default",
			method:     "POST",
			typePrefix: "result",
		},
		"odds": {
			path:       "/football/getJSON.aspx?jsontype=odds_had.aspx",
			method:     "GET",
			typePrefix: "odds",
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			store := testutil.NewMockStore(t, make(map[string][]byte), false)
			ts := testutil.Setup(t, testdataPath)
			defer ts.Close()
			c := &Client{
				Client: client.New(-1),
				Clock:  clockwork.NewFakeClockAt(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			}
			url := fmt.Sprintf("%s%s", ts.URL, tc.path)
			if err := c.Store(context.TODO(), tc.typePrefix, url, tc.method, nil, store); err != nil {
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
			}
		})
	}
}
