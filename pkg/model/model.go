package model

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

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

type Transformer func(io.Reader, Encoder) error

type Modeler struct {
	*bytes.Buffer
	s      store.Store
	itemCh chan Item
}

func New(s store.Store) Modeler {
	buf := new(bytes.Buffer)
	itemCh := make(chan Item)
	return Modeler{buf, s, itemCh}
}

func (m Modeler) Encode(i Model) error {
	var item Item
	if err := i.Item(&item); err != nil {
		return err
	}
	m.itemCh <- item
	return nil
}

func (m Modeler) Update(ctx context.Context, transformers ...Transformer) error {
	errCh := make(chan error)
	wg := sync.WaitGroup{}
	go func() {
		for i := range m.itemCh {
			key := fmt.Sprintf("%s:%s", i.Type, i.Key)
			if err := m.s.Set(key, i.Data); err != nil {
				errCh <- err
			}
		}
		errCh <- nil
	}()
	go func() {
		wg.Wait()
		close(m.itemCh)
	}()
	buffers := make([]io.Writer, len(transformers))
	for i, transform := range transformers {
		buf := new(bytes.Buffer)
		buffers[i] = buf
		wg.Add(1)
		go func(fn Transformer, r io.Reader, encoder Encoder) {
			defer wg.Done()
			if err := fn(r, encoder); err != nil {
				errCh <- err
			}
		}(transform, buf, m)
	}
	w := io.MultiWriter(buffers...)
	go func() {
		if _, err := io.Copy(w, m); err != nil {
			errCh <- err
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	}
	return nil
}
