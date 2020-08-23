package jockey

import (
	"context"
	"fmt"

	"github.com/twistedogic/doom/pkg/store"
)

const (
	typeKey  = "jsontype"
	oddQuery = "odds_%s.aspx"
)

func getOddURL(base, bet string) string {
	terms := make(map[string]string)
	query := fmt.Sprintf(oddQuery, bet)
	terms[typeKey] = query
	return fmt.Sprintf("%s?%s", base, toQueryString(terms))
}

type OddTap struct {
	Client
	bet string
}

func (o OddTap) Update(ctx context.Context, s store.Store) error {
	url := getOddURL(o.BaseURL, o.bet)
	return o.Store(ctx, o.bet, url, "GET", nil, s)
}
