package match

import (
	"encoding/json"
	"io"
)

type MatchModel struct {
	ID      int
	Home    string
	Away    string
	Periods *Periods
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
							m := MatchModel{
								ID:      match.ID,
								Home:    match.Teams.Home.Name,
								Away:    match.Teams.Away.Name,
								Periods: match.Periods,
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
