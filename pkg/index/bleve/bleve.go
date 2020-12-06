package bleve

import (
	"time"
  "strings"

	"github.com/twistedogic/doom/pkg/index"
	"github.com/twistedogic/doom/pkg/store"
	"github.com/twistedogic/doom/proto/model"

	"github.com/blevesearch/bleve"
  "github.com/blevesearch/bleve/search/query"
  "github.com/golang/protobuf/proto"
)

type Meta struct {
	Home, Away string
	Date       time.Time
}

func newMeta(m *model.Match) Meta {
	home := strings.ToLower(m.GetHome().GetName())
	away := strings.ToLower(m.GetAway().GetName())
	ts := m.GetMatchDate().GetSeconds()
	date := time.Unix(ts, 0)
	return Meta{
		Home: home,
		Away: away,
		Date: date,
	}
}

type Index struct {
	store store.Setter
	index bleve.Index
}

func New(path string, target store.Setter) (Index, error) {
	bleveIdx, err := bleve.Open(path)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		bleveIdx, err = bleve.New(path, mapping)
		if err != nil {
			return Index{}, err
		}
	}
	return Index{
		store: target,
		index: bleveIdx,
	}, nil
}

func (i Index) Set(key string, b []byte) error {
	if err := i.store.Set(key, b); err != nil {
		return err
	}
  match := new(model.Match)
  if err := proto.Unmarshal(b, match); err != nil {
    return err
  }
  return i.index.Index(match.GetId(), newMeta(match))
}

func (i Index) search(req *bleve.SearchRequest) ([]string, error) {
  searchResult, err := i.index.Search(req)
  if err != nil {
    return nil, err
  }
  ids := make([]string, 0, len(searchResult.Hits))
  for _, hit := range searchResult.Hits {
    ids = append(ids, hit.ID)
  }
	return ids, nil
}

func (i Index) buildQuery(q index.Query) query.Query {
  queries := []query.Query{}
  if len(q.Home) != 0 {
    home := bleve.NewFuzzyQuery(strings.ToLower(q.Home))
    home.SetField("Home")
    queries = append(queries, home)
  }
  if len(q.Away) != 0 {
    away := bleve.NewFuzzyQuery(strings.ToLower(q.Away))
    away.SetField("Away")
    queries = append(queries, away)
  }
  if !q.Start.IsZero() && !q.End.IsZero() {
    date := bleve.NewDateRangeQuery(q.Start, q.End)
    date.SetField("Date")
    queries = append(queries, date)
  }
  return bleve.NewConjunctionQuery(queries...)
}

func (i Index) Search(q index.Query) ([]string, error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}
  search := bleve.NewSearchRequest(i.buildQuery(q))
  return i.search(search)
}
