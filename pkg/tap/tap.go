package tap

import (
	"context"
)

type Target interface {
	Write([]byte) error
}

type Tap interface {
	Update(context.Context, Target) error
}
