package noop

import (
	hash "crypto/sha1"
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
	h := hash.New()
	h.Reset()
	sum := h.Sum(b)
	key := base64.URLEncoding.EncodeToString(sum[:10])
	return encoder.Encode(Raw{Key: key, Data: b})
}
