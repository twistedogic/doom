package tap

import (
	"context"
	"io"
)

type Tap interface {
	Update(context.Context, io.Writer) error
}
