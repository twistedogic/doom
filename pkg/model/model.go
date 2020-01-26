package model

import (
	"fmt"

	"github.com/twistedogic/doom/pkg/store"
)

type Type string

type Item struct {
	Key  string
	Type Type
	Data []byte
}

type Model interface {
	Item(*Item) error
}

type Encoder interface {
	Encode(Model) error
}

type TransformFunc func([]byte, Encoder) error

type Modeler struct {
	store        store.Store
	transformers []TransformFunc
}

func New(s store.Store, transformers ...TransformFunc) Modeler {
	return Modeler{s, transformers}
}

func (m Modeler) Encode(i Model) error {
	var item Item
	if err := i.Item(&item); err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%s", item.Type, item.Key)
	return m.store.Set(key, item.Data)
}

func (m Modeler) Write(b []byte) error {
	tranformed := false
	for _, fn := range m.transformers {
		if err := fn(b, m); err == nil {
			tranformed = true
		}
	}
	if !tranformed {
		return fmt.Errorf("no transform done")
	}
	return nil
}
