package jockey

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/jonboulle/clockwork"
	"github.com/twistedogic/doom/pkg/client"
	"github.com/twistedogic/doom/pkg/store"
	pb "github.com/twistedogic/doom/proto/source/jockey"
)

const (
	Prefix = "jc"

	baseURL = "https://bet.hkjc.com/football/getJSON.aspx"
)

func toQueryString(kv map[string]string) string {
	terms := make([]string, 0, len(kv))
	for k, v := range kv {
		terms = append(terms, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(terms)
	return strings.Join(terms, "&")
}

func storeProto(key string, match *pb.Match, s store.Store) error {
	b, err := match.Marshal()
	if err != nil {
		return err
	}
	return s.Set(key, b)
}

type Client struct {
	BaseURL string
	client.Client
	clockwork.Clock
}

func New(u string, rate int) Client {
	c := client.New(rate)
	return Client{
		BaseURL: u,
		Client:  c,
		Clock:   clockwork.NewRealClock(),
	}
}

func (c Client) storeProto(typePrefix string, b []byte, s store.Store) error {
	var data []*pb.Data
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for _, matches := range data {
		for _, match := range matches.GetMatches() {
			key := fmt.Sprintf("%s_%s_%s", Prefix, typePrefix, match.GetId())
			if err := storeProto(key, match, s); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c Client) storeRaw(typePrefix string, b []byte, s store.Store) error {
	key := fmt.Sprintf("%s_%s_%d_raw", Prefix, typePrefix, c.Now().Unix())
	return s.Set(key, b)
}

func (c Client) Store(ctx context.Context, typePrefix, url, method string, body io.Reader, s store.Store) error {
	buf := new(bytes.Buffer)
	if err := c.Request(ctx, method, url, body, buf); err != nil {
		return err
	}
	b := buf.Bytes()
	if err := c.storeRaw(typePrefix, b, s); err != nil {
		return err
	}
	return c.storeProto(typePrefix, b, s)
}
