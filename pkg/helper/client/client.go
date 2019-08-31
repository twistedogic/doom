package client

import (
	"io/ioutil"
	"log"
	"net/http"

	json "github.com/json-iterator/go"
	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/helper/flatten"
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
	values := flatten.FlattenDeep(value)
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

type Client struct {
	Limiter ratelimit.Limiter
	BaseURL string
	Rate    int
}

func New(u string, rate int) *Client {
	return &Client{
		BaseURL: u,
		Rate:    rate,
		Limiter: NewLimiter(rate),
	}
}

func (c *Client) updateConfig() {
	limiter := NewLimiter(c.Rate)
	c.Limiter = limiter
}

func (c *Client) Load(s config.Setting) error {
	if err := s.ParseConfig(c); err != nil {
		return err
	}
	c.updateConfig()
	return nil
}

func (c *Client) GetJSON(u string, i interface{}) error {
	c.Limiter.Take()
	return GetJSON(u, i)
}
