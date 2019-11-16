package model

import (
	"context"
	"encoding/json"
	"io"

	"github.com/twistedogic/doom/pkg/store"
)

type Prefix string

type Item struct {
	Type Prefix
	Data json.RawMessage
}

type Modeler interface {
	io.WriteCloser
	Update(context.Context, store.Store) error
}
