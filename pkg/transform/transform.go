package transform

import (
	"github.com/twistedogic/doom/pkg/store"
	"github.com/twistedogic/doom/proto/model"
)

type protoMessage interface {
	Marshal() ([]byte, error)
}

func storeProto(key string, message protoMessage, s store.Setter) error {
	b, err := message.Marshal()
	if err != nil {
		return err
	}
	return s.Set(key, b)
}

type TransformFunc func([]byte, *model.Match) error

type Transformer struct {
	fn TransformFunc
	s  store.Setter
}

func New(fn TransformFunc, s store.Setter) Transformer {
	return Transformer{fn, s}
}

func (t Transformer) storeMatch(match *model.Match) error {
	return storeProto(match.GetId(), match, t.s)
}

func (t Transformer) Set(key string, b []byte) error {
	if err := t.s.Set(key, b); err != nil {
		return err
	}
	match := new(model.Match)
	if err := t.fn(b, match); err != nil {
		return err
	}
	return t.storeMatch(match)
}
