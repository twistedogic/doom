package jockey

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/twistedogic/doom/pkg/client"
	"github.com/twistedogic/doom/pkg/store"
)

const (
	Prefix = "jc"

	baseURL    = "https://bet.hkjc.com/football/getJSON.aspx"
	dateFormat = "2006-01-02"
)

func toQueryString(kv map[string]string) string {
	terms := make([]string, 0, len(kv))
	for k, v := range kv {
		terms = append(terms, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(terms)
	return strings.Join(terms, "&")
}

type container struct {
	Matches []json.RawMessage
}

type Client struct {
	BaseURL string
	client.Client
}

func New(u string, rate int) Client {
	c := client.New(rate)
	return Client{
		BaseURL: u,
		Client:  c,
	}
}

func (c Client) Store(typePrefix string, b []byte, s store.Setter) error {
	h := sha1.New()
	h.Write(b)
	hashKey := hex.EncodeToString(h.Sum(nil)[:12])
	ts := time.Now().Format(dateFormat)
	key := fmt.Sprintf("%s_%s_%s_%s", Prefix, typePrefix, ts, hashKey)
	return s.Set(key, b)
}

func (c Client) StoreMatch(ctx context.Context, typePrefix, url, method string, body io.Reader, s store.Setter) error {
	buf := new(bytes.Buffer)
	if err := c.Request(ctx, method, url, body, buf); err != nil {
		return err
	}
	var containers []container
	if err := json.Unmarshal(buf.Bytes(), &containers); err != nil {
		return err
	}
	for _, data := range containers {
		for _, m := range data.Matches {
			if err := c.Store(typePrefix, []byte(m), s); err != nil {
				return err
			}
		}
	}
	return nil
}
