package csv

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/helper/file"
	"github.com/twistedogic/doom/pkg/helper/flatten"
)

type Target struct {
	Path      string
	writer    *csv.Writer
	info      os.FileInfo
	hasHeader bool
}

func New(filename string) (*Target, error) {
	f, err := file.CreateIfNotExist(filename)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &Target{
		filename,
		csv.NewWriter(f),
		info,
		info.Size() != 0,
	}, err
}

func (t *Target) Load(s config.Setting) error {
	if err := s.ParseConfig(t); err != nil {
		return err
	}
	f, err := file.CreateIfNotExist(t.Path)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		return err
	}
	t.writer = csv.NewWriter(f)
	t.info = info
	t.hasHeader = info.Size() != 0
	return err
}

func (t *Target) UpsertItem(i interface{}) error {
	if !t.hasHeader {
		t.hasHeader = true
		return Marshal(t.writer, i, true)
	}
	return Marshal(t.writer, i, false)
}

func (t *Target) BulkUpsert(i interface{}) error {
	items, ok := flatten.InterfaceToSlice(i)
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
