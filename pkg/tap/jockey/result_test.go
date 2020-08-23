package jockey

import (
	"reflect"
	"testing"
	"time"
)

func Test_chunkTime(t *testing.T) {
	cases := map[string]struct {
		start, end time.Time
		step       time.Duration
		want       [][]time.Time
	}{
		"base": {
			start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			step:  5 * 24 * time.Hour,
			want: [][]time.Time{
				[]time.Time{
					time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				},
				[]time.Time{
					time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
					time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			got := chunkTime(tc.start, tc.end, tc.step)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
