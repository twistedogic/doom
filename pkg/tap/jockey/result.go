package jockey

import (
	"context"
	"fmt"
	"time"

	"github.com/twistedogic/doom/pkg/store"
)

const (
	resultDateFormat = "20060102"

	resultQuery  = "search_result.aspx"
	startDateKey = "startdate"
	endDateKey   = "enddate"

	MONTH = 28 * 24 * time.Hour
)

func getResultURL(base string, start, end time.Time) string {
	terms := make(map[string]string)
	terms[typeKey] = resultQuery
	terms[startDateKey] = start.Format(resultDateFormat)
	terms[endDateKey] = end.Format(resultDateFormat)
	return fmt.Sprintf("%s?%s&teamid=default", base, toQueryString(terms))
}

func chunkTime(start, end time.Time, step time.Duration) [][]time.Time {
	diff := end.Sub(start) / step
	out := make([][]time.Time, 0, diff+1)
	current := start
	for {
		s := current
		e := current.Add(step)
		if e.After(end) {
			e := end
			out = append(out, []time.Time{s, e})
			break
		}
		out = append(out, []time.Time{s, e})
		current = e
	}
	return out
}

type ResultTap struct {
	Client
	start, end time.Time
}

func NewResultTap(base string, rate int, start, end time.Time) ResultTap {
	c := New(base, rate)
	return ResultTap{c, start, end}
}

func (r ResultTap) Update(ctx context.Context, s store.Setter) error {
	chunks := chunkTime(r.start, r.end, MONTH)
	for _, chunk := range chunks {
		url := getResultURL(r.BaseURL, chunk[0], chunk[1])
		if err := r.StoreMatch(ctx, "result", url, "POST", nil, s); err != nil {
			return err
		}
	}
	return nil
}
