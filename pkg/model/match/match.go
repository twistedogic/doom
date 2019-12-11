package match

import (
	"encoding/json"
	"io"
	"time"
)

type MatchModel struct {
	ID          int
	Home        string
	Away        string
	Periods     *Periods
	Result      Result
	LastUpdated time.Time
}

func Transform(r io.Reader, target io.WriteCloser) error {
	defer target.Close()
	decoder := json.NewDecoder(r)
	encoder := json.NewEncoder(target)
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
							m := MatchModel{
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
