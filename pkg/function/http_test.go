package function

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOddHandler(t *testing.T) {
	handler := OddHTTP
	req := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("%s", body)
	}
	t.Logf("%s", body)
}
