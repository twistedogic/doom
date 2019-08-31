package tap

import (
	"encoding/json"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/fatih/structs"
	"github.com/twistedogic/doom/pkg/config"
	"github.com/twistedogic/doom/pkg/target"
)

type Tap interface {
	Load(config.Setting) error
	Update(target.Target) error
}

type MessageType string

const (
	SCHEMA MessageType = "SCHEMA"
	RECORD MessageType = "SCHEMA"
	STATE  MessageType = "STATE"
)

type Message struct {
	Type   MessageType     `json:"type"`
	Stream string          `json:"stream,omitempty"`
	Schema json.RawMessage `json:"schema,omitemtpy"`
	Record json.RawMessage `json:"record,omitemtpy"`
	Value  json.RawMessage `json:"value,omitemtpy"`
}

func NewSchema(i interface{}, m *Message) error {
	m.Type = SCHEMA
	m.Stream = strings.ToLower(structs.Name(i))
	schema := jsonschema.Reflect(i)
	b, err := json.Marshal(schema.Definitions)
	if err != nil {
		return err
	}
	m.Schema = json.RawMessage(b)
	return nil

}
func NewRecord(i interface{}, m *Message) error {
	m.Type = RECORD
	m.Stream = strings.ToLower(structs.Name(i))
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	m.Record = json.RawMessage(b)
	return nil
}

func NewState(i interface{}, m *Message) error {
	m.Type = STATE
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	m.Value = json.RawMessage(b)
	return nil
}
