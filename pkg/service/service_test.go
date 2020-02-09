package service

import (
	"testing"

	"github.com/twistedogic/doom/proto/env"
)

func TestHelloWorld(t *testing.T) {
	t.Log(env.Name{Value: "test"})
}
