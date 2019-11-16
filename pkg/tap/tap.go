package tap

import (
	"context"
	"io"
)

type Tap interface {
	Update(context.Context, io.WriteCloser) error
}
