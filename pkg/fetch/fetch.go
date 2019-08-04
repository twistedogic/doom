package fetch

import (
	"encoding/json"

	"github.com/twistedogic/doom/pkg/helper"
	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/pkg/service/radar"
	"github.com/twistedogic/jsonpath"
	"go.uber.org/ratelimit"
)

const (
	MatchPath  = "$.doc[*].data[*].realcategories[*].tournaments[*].matches"
	DetailPath = "$.doc[*].data[*]"
)

type Fetcher struct {
	*radar.Client
	ratelimit.Limiter
}

func New(u string, rate int) *Fetcher {
	c := radar.New(u)
	limiter := ratelimit.New(rate)
	if rate == -1 {
		limiter = ratelimit.NewUnlimited()
	}
	return &Fetcher{c, limiter}
}

func (f *Fetcher) ExtractJsonPath(i interface{}, path string) ([][]byte, error) {
	value, err := jsonpath.Lookup(path, i)
	if err != nil {
		return nil, err
	}
	values := helper.FlattenDeep(value)
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

func (f *Fetcher) GetMatch(offset int) ([]model.Match, error) {
	f.Take()
	var container interface{}
	if err := f.GetMatchFullFeed(offset, &container); err != nil {
		return nil, err
	}
	items, err := f.ExtractJsonPath(container, MatchPath)
	if err != nil {
		return nil, err
	}
	out := make([]model.Match, len(items))
	for i, item := range items {
		if err := jsonpath.Unmarshal(item, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (f *Fetcher) GetDetail(matchID int) ([]model.Detail, error) {
	f.Take()
	var container interface{}
	if err := f.GetMatchDetail(matchID, &container); err != nil {
		return nil, err
	}
	items, err := f.ExtractJsonPath(container, DetailPath)
	if err != nil {
		return nil, err
	}
	out := make([]model.Detail, len(items))
	for i, item := range items {
		if err := jsonpath.Unmarshal(item, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}
