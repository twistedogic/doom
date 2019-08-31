package config

import (
	"io"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Tap      Setting `yaml:"tap"`
	Target   Setting `yaml:"target"`
	Schedule string  `yaml:"schedule"`
}

type Setting struct {
	Name   string            `yaml:"name"`
	Config map[string]string `yaml:"config"`
}

func (s Setting) ParseConfig(i interface{}) error {
	return mapstructure.WeakDecode(s.Config, i)
}

func New(b []byte) (Config, error) {
	var c Config
	err := yaml.UnmarshalStrict(b, &c)
	return c, err
}

func Load(r io.Reader) (Config, error) {
	var c Config
	decoder := yaml.NewDecoder(r)
	decoder.SetStrict(true)
	err := decoder.Decode(&c)
	return c, err
}
