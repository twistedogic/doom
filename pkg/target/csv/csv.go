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

func marshal(w *csv.Writer, rows []Row) error {
	if len(rows) == 0 {
		return fmt.Errorf("rows is empty")
	}
	header := rows[0].Header
	if err := w.Write(header); err != nil {
		return err
	}
	for _, r := range rows {
		if err := w.Write(r.Values); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func Marshal(w *csv.Writer, i interface{}) error {
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
	return marshal(w, rows)
}
