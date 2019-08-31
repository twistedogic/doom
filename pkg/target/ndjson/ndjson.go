package ndjson

import (
	"bufio"
	"io"
	"os"
	"time"

	json "github.com/json-iterator/go"
	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/helper/file"
	"github.com/twistedogic/doom/pkg/helper/flatten"
	"github.com/twistedogic/doom/pkg/target"
)

type Target struct {
	Path string
	file *os.File
}

func New(filename string) (*Target, error) {
	f, err := file.CreateIfNotExist(filename)
	if err != nil {
		return nil, err
	}
	return &Target{filename, f}, nil
}

func (t *Target) Load(s config.Setting) error {
	if err := s.ParseConfig(t); err != nil {
		return err
	}
	f, err := file.CreateIfNotExist(t.Path)
	if err != nil {
		return err
	}
	t.file = f
	return nil
}

func (t *Target) UpsertItem(i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	if _, err := t.file.Write(b); err != nil {
		return err
	}
	if _, err := t.file.WriteString("\n"); err != nil {
		return err
	}
	return nil
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
	info, _ := t.file.Stat()
	return info.ModTime()
}

func (t *Target) Update(dst target.Target) error {
	r := bufio.NewReader(t.file)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(line) != 0 {
			var i interface{}
			if err := json.Unmarshal(line, &i); err != nil {
				return err
			}
			if err := dst.UpsertItem(i); err != nil {
				return err
			}
		}
	}
	return nil
}
