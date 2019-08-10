package csv

import (
	"encoding/csv"
	"fmt"
	"sort"

	"github.com/twistedogic/doom/pkg/helper"
)

var DefaultKeyDelimiter string = "_"

type Marshaler interface {
	MarshalCSV() map[string]string
}

type Row struct {
	Header []string
	Values []string
}

func newRow(m map[string]string) Row {
	header := make([]string, 0, len(m))
	values := make([]string, 0, len(m))
	for k := range m {
		header = append(header, k)
	}
	sort.Strings(header)
	for _, k := range header {
		values = append(values, m[k])
	}
	return Row{header, values}
}

func marshalHeader(w *csv.Writer, row Row) error {
	return w.Write(row.Header)
}

func marshal(w *csv.Writer, rows []Row) error {
	if len(rows) == 0 {
		return fmt.Errorf("rows is empty")
	}
	for _, r := range rows {
		if err := w.Write(r.Values); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func Marshal(w *csv.Writer, i interface{}, withHeader bool) error {
	rowMaps := []map[string]string{}
	if val, ok := helper.InterfaceToSlice(i); ok {
		for _, v := range val {
			switch t := v.(type) {
			case Marshaler:
				rowMaps = append(rowMaps, t.MarshalCSV())
			default:
				rowMaps = append(rowMaps, helper.FlattenKey(t, DefaultKeyDelimiter)...)
			}
		}
	} else {
		rowMaps = helper.FlattenKey(i, DefaultKeyDelimiter)
	}
	rows := make([]Row, len(rowMaps))
	for j, m := range rowMaps {
		rows[j] = newRow(m)
	}
	if withHeader {
		if err := marshalHeader(w, rows[0]); err != nil {
			return err
		}
	}
	return marshal(w, rows)
}
