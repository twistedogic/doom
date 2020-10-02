package jockey

import (
	"testing"

	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/testutil"
)

func setupOddTapTest(t *testing.T, url string) tap.Tap {
	t.Helper()
	return NewOddTap(url, "had", -1)
}

func Test_OddTap_Update(t *testing.T) {
	testutil.TapTest(t, testdataPath, setupOddTapTest, false)
}
