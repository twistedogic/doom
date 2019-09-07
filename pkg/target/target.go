package target

import (
	"github.com/twistedogic/doom/pkg/config"
)

type Target interface {
	Load(config.Setting) error
	UpsertItem(interface{}) error
	Close() error
}
