package csv

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/twistedogic/doom/pkg/helper"
)

type Target struct {
	writer *csv.Writer
	info   os.FileInfo
}

func New(filename string) (*Target, error) {
	info, err := os.Stat(filename)
	if err == nil {
		if f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return nil, err
		} else {
			return &Target{csv.NewWriter(f), info}, nil
		}
	} else if os.IsNotExist(err) {
		if f, err := os.Create(filename); err != nil {
			return nil, err
		} else {
			info, _ := f.Stat()
			return &Target{csv.NewWriter(f), info}, nil
		}
	}
	return nil, err
}

func (t *Target) UpsertItem(i interface{}) error {
	return Marshal(t.writer, i)
}

func (t *Target) BulkUpsert(i interface{}) error {
	items, ok := helper.InterfaceToSlice(i)
	if !ok {
		return t.UpsertItem(i)
	}
	for _, item := range items {
		if err := t.UpsertItem(item); err != nil {
			return err
		}
	}
	return nil
}

func (t *Target) GetLastUpdate() time.Time {
	return t.info.ModTime()
}
