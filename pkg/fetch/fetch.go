package fetch

import (
	"encoding/json"

	"github.com/oliveagle/jsonpath"
	"github.com/twistedogic/doom/pkg/schema"
	"github.com/twistedogic/doom/pkg/schema/schemautil"
	"github.com/twistedogic/doom/pkg/service/radar"
)

const (
	MatchPath  = "$.doc[*].data[*].realcategories[*].tournaments[*].matches"
	DetailPath = "$.doc[*].data[*]"
)

type Fetcher struct {
	*radar.Client
}

func New(u string) *Fetcher {
	c := radar.New(u)
	return &Fetcher{c}
}

func (f *Fetcher) ExtractJsonPath(i interface{}, path string) ([][]byte, error) {
	value, err := jsonpath.JsonPathLookup(i, path)
	if err != nil {
		return nil, err
	}
	values := schemautil.FlattenDeep(value)
	out := make([][]byte, len(values))
	for i, v := range values {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		out[i] = b
	}
	return out, nil
}

func (f *Fetcher) GetMatch(offset int) ([]schema.Match, error) {
	var container interface{}
	if err := f.GetMatchFullFeed(offset, &container); err != nil {
		return nil, err
	}
	items, err := f.ExtractJsonPath(container, MatchPath)
	if err != nil {
		return nil, err
	}
	out := make([]schema.Match, len(items))
	for i, item := range items {
		if err := json.Unmarshal(item, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (f *Fetcher) GetDetail(matchID uint64) ([]schema.Detail, error) {
	var container interface{}
	if err := f.GetMatchDetail(int(matchID), &container); err != nil {
		return nil, err
	}
	items, err := f.ExtractJsonPath(container, DetailPath)
	if err != nil {
		return nil, err
	}
	out := make([]schema.Detail, len(items))
	for i, item := range items {
		if err := json.Unmarshal(item, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}
