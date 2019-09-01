package model

import (
	"encoding/csv"
	"io"
	"strings"

	json "github.com/json-iterator/go"
)

type Entry struct {
	Header []string
	Values []string
}

func New(header, values string) Entry {
	return Entry{
		Header: strings.Split(header, ","),
		Values: strings.Split(values, ","),
	}
}

func (e Entry) MarshalCSV() map[string]string {
	out := make(map[string]string)
	for i, key := range e.Header {
		out[key] = e.Values[i]
	}
	return out
}

func (e Entry) MarshalJSON() ([]byte, error) {
	m := e.MarshalCSV()
	return json.Marshal(m)
}

type Decoder struct {
	r      io.Reader
	header []string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode(ch chan Entry) error {
	defer close(ch)
	r := csv.NewReader(d.r)
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(line) == 0 {
			continue
		}
		if len(d.header) == 0 {
			d.header = line
			continue
		}
		ch <- Entry{d.header, line}
	}
	return nil
}
