package model

import (
	"context"
	"encoding/json"
	"io"

	"github.com/twistedogic/doom/pkg/store"
)

type Type string

type Item interface {
	Key() string
	Type() Type
	Data() []byte
}

type Modeler interface {
	io.WriteCloser
	Update(context.Context, store.Store) error
}
