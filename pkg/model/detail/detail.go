package detail

import (
	"strconv"

	json "github.com/json-iterator/go"
	"github.com/twistedogic/doom/pkg/model"
)

const Type model.Type = "detail"

type Model struct {
	MatchID int
	Teams   Teams
	Values  []Value
}

func (m Model) Item(i *model.Item) error {
	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	i.Key, i.Type, i.Data = strconv.Itoa(m.MatchID), Type, b
	return nil
}

func Transform(b []byte, encoder model.Encoder) error {
	var feed Detail
	if err := json.Unmarshal(b, &feed); err != nil {
		return err
	}
	for _, doc := range feed.Doc {
		values := make([]Value, 0, len(doc.Data.Values))
		for _, v := range doc.Data.Values {
			values = append(values, v)
		}
		model := Model{
			MatchID: doc.Data.MatchID,
			Teams:   doc.Data.Teams,
			Values:  values,
		}
		if err := encoder.Encode(model); err != nil {
			return err
		}
	}
	return nil
}
