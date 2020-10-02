package tap

import (
	"context"

	"github.com/twistedogic/doom/pkg/store"
)

type Tap interface {
	Update(context.Context, store.Setter) error
}

type TapOperation struct {
	Tap
	store.Store
}

func (t TapOperation) Run(ctx context.Context) error {
	return t.Update(ctx, t)
}
