package detail

import (
	"encoding/json"
	"io"
)

type DetailModel struct {
	MatchID int
	Teams   Teams
	Values  []Value
}

func Transform(r io.Reader, target io.WriteCloser) error {
	defer target.Close()
	decoder := json.NewDecoder(r)
	encoder := json.NewEncoder(target)
	for decoder.More() {
		var feed Detail
		if err := decoder.Decode(&feed); err != nil {
			return err
		}
		for _, doc := range feed.Doc {
			values := make([]Value, 0, len(doc.Data.Values))
			for _, v := range doc.Data.Values {
				values = append(values, v)
			}
			model := DetailModel{
				MatchID: doc.Data.MatchID,
				Teams:   doc.Data.Teams,
				Values:  values,
			}
			if err := encoder.Encode(model); err != nil {
				return err
			}
		}
	}
	return nil
}
