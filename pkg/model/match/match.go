package match

import (
	"io"
	"strconv"
	"time"

	json "github.com/json-iterator/go"
	"github.com/twistedogic/doom/pkg/model"
)

const Type model.Type = "match"

type Model struct {
	ID          int
	Home        string
	Away        string
	Periods     *Periods
	Result      Result
	LastUpdated time.Time
}

func (m Model) Item(i *model.Item) error {
	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	key := strconv.Itoa(m.ID)
	i.Key, i.Type, i.Data = key, Type, b
	return nil
}

func Transform(r io.Reader, encoder model.Encoder) error {
	decoder := json.NewDecoder(r)
	for decoder.More() {
		var feed Feed
		if err := decoder.Decode(&feed); err != nil {
			return err
		}
		for _, doc := range feed.Doc {
			for _, datum := range doc.Data {
				for _, cat := range datum.Realcategories {
					for _, tournament := range cat.Tournaments {
						for _, match := range tournament.Matches {
							lastUpdate := time.Unix(int64(match.UpdatedUts), 0)
							m := Model{
								ID:          match.ID,
								Home:        match.Teams.Home.Name,
								Away:        match.Teams.Away.Name,
								Periods:     match.Periods,
								Result:      match.Result,
								LastUpdated: lastUpdate,
							}
							if err := encoder.Encode(m); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}
