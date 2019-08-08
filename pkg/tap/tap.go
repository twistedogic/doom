package tap

import (
	"time"

	"github.com/twistedogic/doom/pkg/target"
)

type Tap interface {
	Update(time.Time, target.Target) error
}
