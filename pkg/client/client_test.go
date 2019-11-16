package client

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Setup(t *testing.T, body []byte) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write(body)
	}))
}

func TestRequest(t *testing.T) {
	payload := []byte(`something`)
	ts := Setup(t, payload)
	defer ts.Close()
	buf := &bytes.Buffer{}
	client := New(0)
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Request(req, buf); err != nil {
		t.Fatal(err)
	}
	if string(payload) != string(buf.Bytes()) {
		t.Fail()
	}
}

func TestGetResponse(t *testing.T) {
	payload := []byte(`something`)
	ts := Setup(t, payload)
	defer ts.Close()
	buf := &bytes.Buffer{}
	client := New(0)
	if err := client.GetResponse(ts.URL, buf); err != nil {
		t.Fatal(err)
	}
	if string(payload) != string(buf.Bytes()) {
		t.Fail()
	}
}
