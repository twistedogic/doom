package radar

import (
	"testing"
)

func TestClient(t *testing.T) {
	client := New(RadarURL)
	var out interface{}
	if err := client.GetMatchFullFeed(0, &out); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", out)
}
