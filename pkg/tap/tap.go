package tap

import (
	"github.com/twistedogic/doom/pkg/target"
)

type Tap interface {
	Update(target.Target) error
}
