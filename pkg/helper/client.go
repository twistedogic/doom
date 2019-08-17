package helper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/twistedogic/jsonpath"
	"go.uber.org/ratelimit"
)

func NewLimiter(rate int) ratelimit.Limiter {
	if rate > 0 {
		return ratelimit.New(rate)
	}
	return ratelimit.NewUnlimited()
}

func GetJSON(u string, value interface{}) error {
	res, err := http.Get(u)
	if err != nil {
		log.Printf("ERROR GET %s %s", u, err)
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("ERROR GET %s %s", u, err)
		return err
	}
	if err := json.Unmarshal(b, value); err != nil {
		log.Printf("ERROR GET %s %s", u, err)
		return err
	}
	return nil
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
