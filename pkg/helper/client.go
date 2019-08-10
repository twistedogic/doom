package helper

import (
	"encoding/json"
	"net/http"

	"github.com/twistedogic/jsonpath"
)

func GetJSON(u string, value interface{}) error {
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(value)
}

func ExtractJsonPath(i interface{}, path string) ([][]byte, error) {
	value, err := jsonpath.Lookup(path, i)
	if err != nil {
		return nil, err
	}
	values := FlattenDeep(value)
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
