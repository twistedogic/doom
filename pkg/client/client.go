package client

import (
	"encoding/json"
	"net/http"
)

func GetJSON(u string, value interface{}) error {
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(value)
}
