package model

import (
	"context"
	"encoding/json"
	"io"

	"github.com/twistedogic/doom/pkg/store"
)

type Type string

type Item struct {
	Type Type
	Data json.RawMessage
}

type Modeler interface {
	io.WriteCloser
	Update(context.Context, store.Store) error
}
