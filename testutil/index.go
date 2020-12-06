package testutil

import (
	"sort"
	"testing"
	"time"

	"github.com/twistedogic/doom/pkg/index"
	"github.com/twistedogic/doom/proto/model"

	"github.com/gogo/protobuf/types"
	"github.com/google/go-cmp/cmp"
)

func toTimestamp(t time.Time) *types.Timestamp {
	ts, _ := types.TimestampProto(t)
	return ts
}

type searchCase struct {
	query index.Query
	want  []string
}

func (s searchCase) Check(t *testing.T, i index.Index) {
	got, err := i.Search(s.query)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	sort.Strings(s.want)
	if diff := cmp.Diff(s.want, got); diff != "" {
		t.Fatalf("for\n%v\n%s", s.query, diff)
	}
}

func IndexTest(t *testing.T, target index.Index, cleanup func()) {
	defer cleanup()
	input := []*model.Match{
		&model.Match{
			Id:        "a",
			Home:      &model.Team{Name: "something"},
			Away:      &model.Team{Name: "otherthing"},
			MatchDate: toTimestamp(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		&model.Match{
			Id:        "b",
			Home:      &model.Team{Name: "anything"},
			Away:      &model.Team{Name: "otherthing"},
			MatchDate: toTimestamp(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)),
		},
		&model.Match{
			Id:        "c",
			Home:      &model.Team{Name: "nothing"},
			Away:      &model.Team{Name: "otherthing"},
			MatchDate: toTimestamp(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)),
		},
	}
	cases := map[string]searchCase{
		"empty":       {index.Query{}, []string{}},
		"case insenitive": {index.Query{Home: "SomeThing"}, []string{"a"}},
		"exact match": {index.Query{Home: "something"}, []string{"a"}},
		"match": {
			index.Query{
				Away: "otherthing",
			},
			[]string{"a", "b", "c"},
		},

		"time range": {
			index.Query{
				Start: time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
			},
			[]string{"b"},
		},
		"mix": {
			index.Query{
				Away:  "otherthing",
				Start: time.Date(2019, 12, 30, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
			},
			[]string{"a", "b"},
		},
		"no result": {
			index.Query{
				Home:  "something",
				Start: time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
			},
			[]string{},
		},
	}
	for _, entry := range input {
		id := entry.GetId()
		b, err := entry.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if err := target.Set(id, b); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range cases {
		tc.Check(t, target)
	}
}
