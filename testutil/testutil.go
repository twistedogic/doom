package testutil

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

var (
	TestError = errors.New("test error")
)

func Setup(t *testing.T, path string) *httptest.Server {
	t.Helper()
	log.SetOutput(ioutil.Discard)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	data := make(map[string][]byte)
	for _, info := range files {
		b, err := ioutil.ReadFile(filepath.Join(path, info.Name()))
		if err != nil {
			t.Fatal(err)
		}
		data[info.Name()] = b
	}
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		for k, b := range data {
			lower := strings.ToLower(req.URL.Path)
			if strings.Contains(lower, k) {
				res.Write(b)
				return
			}
		}
		http.Error(res, req.URL.Path, 500)
	}))
}

func Write(t *testing.T, filename string, r io.Reader) {
	t.Helper()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filename, b, 0644); err != nil {
		t.Fatal(err)
	}
}

func Open(t *testing.T, dir, name string, repeat int) *bytes.Buffer {
	t.Helper()
	b, err := ioutil.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	for i := 0; i < repeat; i++ {
		if _, err := buf.Write(b); err != nil {
			t.Fatal(err)
		}
	}
	return buf
}

type MockTarget struct {
	*bytes.Buffer
	t        *testing.T
	showLog  bool
	hasError bool
}

func NewMockTarget(t *testing.T, showLog, hasError bool) MockTarget {
	t.Helper()
	return MockTarget{&bytes.Buffer{}, t, showLog, hasError}
}

func (m MockTarget) Close() error {
	if m.showLog {
		m.t.Log(string(m.Bytes()))
	}
	if m.hasError {
		return TestError
	}
	return nil
}

type MockStore struct {
	t        *testing.T
	content  map[string][]byte
	hasError bool
}

func NewMockStore(t *testing.T, content map[string][]byte, hasError bool) *MockStore {
	return &MockStore{t, content, hasError}
}

func (m *MockStore) Content() map[string][]byte {
	return m.content
}

func (m *MockStore) Get(key string) ([]byte, error) {
	return m.content[key], nil
}

func (m *MockStore) Set(key string, b []byte) error {
	if m.hasError {
		return TestError
	}
	m.content[key] = b
	return nil
}

func (m *MockStore) Scan(...string) ([]string, error) {
	if m.hasError {
		return nil, TestError
	}
	keys := make([]string, 0, len(m.content))
	for k := range m.content {
		keys = append(keys, k)
	}
	return keys, nil
}
