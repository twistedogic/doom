package config

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name     string  `yaml:"name"`
	Tap      Setting `yaml:"tap"`
	Target   Setting `yaml:"target"`
	Schedule string  `yaml:"schedule"`
	Timeout  string  `yaml:"timeout"`
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

func Load(filename string) (Config, error) {
	var c Config
	f, err := os.Open(filename)
	if err != nil {
		return c, err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	err = decoder.Decode(&c)
	return c, err
}
