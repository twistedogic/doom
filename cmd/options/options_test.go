package options

import (
	"testing"
)

func TestTapOptions(t *testing.T) {
	if len(TapOptions) == 0 {
		t.Fail()
	}
}

func TestTargetOptions(t *testing.T) {
	if len(TargetOptions) == 0 {
		t.Fail()
	}
}
