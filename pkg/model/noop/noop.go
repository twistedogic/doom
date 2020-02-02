package noop

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/twistedogic/doom/pkg/model"
)

const Type model.Type = "raw"

type Raw struct {
	Key  string
	Data []byte
}

func (r Raw) Item(i *model.Item) error {
	i.Key, i.Type, i.Data = r.Key, Type, r.Data
	return nil
}

func Transform(b []byte, encoder model.Encoder) error {
	hash := sha1.New()
	hash.Reset()
	sum := hash.Sum(b)
	key := base64.URLEncoding.EncodeToString(sum)
	return encoder.Encode(Raw{Key: key, Data: b})
}
