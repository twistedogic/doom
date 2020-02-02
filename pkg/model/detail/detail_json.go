package detail

import (
	"encoding/json"
)

type Detail struct {
	Doc []Doc `json:"doc"`
}

type Doc struct {
	Data Data `json:"data"`
}

type Data struct {
	MatchID int    `json:"_matchid"`
	Teams   Teams  `json:"teams"`
	Values  Values `json:"values"`
}

type Teams struct {
	Away string `json:"away"`
	Home string `json:"home"`
}

type Value struct {
	Name string
	Away int
	Home int
}

type Values []Value

func (vs *Values) UnmarshalJSON(b []byte) error {
	container := make(map[string]struct {
		Name  string
		Value map[string]json.Number
	})
	if err := json.Unmarshal(b, &container); err != nil {
		return err
	}
	values := make([]Value, 0, len(container))
	for _, entry := range container {
		item := Value{Name: entry.Name}
		for k, number := range entry.Value {
			var value int
			if v, err := number.Int64(); err == nil {
				value = int(v)
			}
			switch k {
			case "home":
				item.Home = value
			case "away":
				item.Away = value
			}
		}
		values = append(values, item)
	}
	*vs = values
	return nil
}
