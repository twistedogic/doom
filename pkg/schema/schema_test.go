package schema

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/fatih/structs"
	"github.com/google/go-cmp/cmp"
	"github.com/oliveagle/jsonpath"
	"github.com/twistedogic/doom/pkg/schema/schemautil"
)

const testdataPath = "../../testdata"

func ReadTestdata(t *testing.T, file string) interface{} {
	t.Helper()
	b, err := ioutil.ReadFile(filepath.Join(testdataPath, file))
	if err != nil {
		t.Fatal(err)
	}
	var container interface{}
	if err := json.Unmarshal(b, &container); err != nil {
		t.Fatal(err)
	}
	return container
}

func TestParseMatch(t *testing.T) {
	container := ReadTestdata(t, "fullfeed")
	value, err := jsonpath.JsonPathLookup(container, "$.doc[*].data[*].realcategories[*].tournaments[*].matches")
	if err != nil {
		t.Fatal(err)
	}
	values := schemautil.FlattenDeep(value)
	for _, v := range values {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		var m Match
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatal(err)
		}
		if cmp.Equal(Match{}, m) || structs.HasZero(m) {
			t.Fatalf("not parsed")
		}
	}
}

func TestParseDetail(t *testing.T) {
	container := ReadTestdata(t, "details")
	value, err := jsonpath.JsonPathLookup(container, "$.doc[*].data[*]")
	if err != nil {
		t.Fatal(err)
	}
	values := schemautil.FlattenDeep(value)
	for _, v := range values {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		var d Detail
		if err := json.Unmarshal(b, &d); err != nil {
			t.Fatal(err)
		}
		if cmp.Equal(Detail{}, d) || structs.HasZero(d) {
			t.Fatalf("not parsed")
		}
	}
}
