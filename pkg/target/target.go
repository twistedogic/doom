package target

import (
	"time"

	"github.com/twistedogic/doom/pkg/config"
)

type Target interface {
	Load(config.Setting) error
	UpsertItem(interface{}) error
	BulkUpsert(interface{}) error
	GetLastUpdate() time.Time
}
