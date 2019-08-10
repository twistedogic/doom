package target

import (
	"time"
)

type Target interface {
	UpsertItem(interface{}) error
	BulkUpsert(interface{}) error
	GetLastUpdate() time.Time
}
