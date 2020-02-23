package store

import (
	"reflect"
	"strings"

	json "github.com/json-iterator/go"
)

type (
	Field struct {
		name string
		i    interface{}
	}

	FieldTermPair struct {
		Field string
		Term  string
	}
)

func (f Field) Name() string {
	return f.name
}

func (f Field) Kind() reflect.Kind {
	return reflect.TypeOf(f.i).Kind()
}

func (f Field) Value() interface{} {
	return f.i
}

func GetFields(i interface{}) []Field {
	if reflect.TypeOf(i).Kind() != reflect.Map {
		return make([]Field, 0)
	}
	m, ok := i.(map[string]interface{})
	if !ok {
		return make([]Field, 0)
	}
	fields := make([]Field, 0, len(m))
	for k, v := range m {
		fields = append(fields, Field{k, v})
	}
	return fields
}

func ParseFieldTermPairs(b []byte) ([]FieldTermPair, error) {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return nil, err
	}
	output := make([]FieldTermPair, 0)
	var field Field
	fields := GetFields(item)
	for len(fields) != 0 {
		field, fields = fields[0], fields[1:]
		fieldName := strings.ToLower(field.Name())
		fieldValue := field.Value()
		switch field.Kind() {
		case reflect.String:
			value := fieldValue.(string)
			output = append(output, FieldTermPair{
				Field: fieldName,
				Term:  strings.ToLower(value),
			})
		case reflect.Map:
			nestedFields := GetFields(fieldValue)
			fields = append(fields, nestedFields...)
		case reflect.Slice:
			values := fieldValue.([]interface{})
			for _, value := range values {
				nestedField := Field{fieldName, value}
				fields = append(fields, nestedField)
			}
		}
	}
	return output, nil
}
