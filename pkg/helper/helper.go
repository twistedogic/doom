package helper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

func IsIterable(v interface{}) bool {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Slice:
		return true
	case reflect.Array:
		return true
	default:
		return false
	}
}

func FlattenDeep(args ...interface{}) []interface{} {
	list := []interface{}{}
	for _, v := range args {
		if IsIterable(v) {
			for _, z := range FlattenDeep((v.([]interface{}))...) {
				list = append(list, z)
			}
		} else {
			list = append(list, v)
		}
	}
	return list
}

func SetField(obj interface{}, name string, value interface{}) error {
	field, ok := structs.New(obj).FieldOk(name)
	if !ok {
		return fmt.Errorf("no field %s found in %#v", name, obj)
	}
	return field.Set(value)
}

type Field struct {
	Tag         string
	Name        string
	Value       interface{}
	IsZeroValue bool
}

func GetTaggedFields(s interface{}, tagName string) []Field {
	sfields := structs.New(s).Fields()
	fields := make([]Field, 0, len(sfields))
	for _, f := range sfields {
		tag := f.Tag(tagName)
		if tag != "" {
			fields = append(fields, Field{
				Tag:         tag,
				Name:        f.Name(),
				Value:       f.Value(),
				IsZeroValue: f.IsZero(),
			})
		}
	}
	return fields
}

func GetFieldName(v interface{}, tagName, tagValue string) (string, error) {
	//TODO: for finding embedded primary key
	for _, field := range GetTaggedFields(v, tagName) {
		if strings.Contains(field.Tag, tagValue) {
			return field.Name, nil
		}
	}
	return "", fmt.Errorf("No field found with tag name: %s and tag value: %s", tagName, tagValue)
}
