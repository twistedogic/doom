package schemautil

import (
	"fmt"
	"reflect"
	"strings"
)

type Field struct {
	Tag         string
	Name        string
	Value       reflect.Value
	IsZeroValue bool
}

func isPtr(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Ptr
}

func GetTaggedFields(s interface{}, tagName string) []Field {
	fields := make([]Field, 0)
	ifValue := reflect.ValueOf(s)
	ifType := reflect.TypeOf(s)
	if isPtr(s) {
		ifValue = reflect.Indirect(ifValue)
		ifType = ifValue.Type()
	}
	for i := 0; i < ifType.NumField(); i++ {
		t := ifType.Field(i)
		v := ifValue.Field(i)
		isValid := v.IsValid() && v.CanInterface()
		isZeroValue := reflect.Zero(v.Type()).Interface() == v.Interface()
		if tag, ok := t.Tag.Lookup(tagName); ok && isValid {
			fields = append(fields, Field{tag, t.Name, v, isZeroValue})
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
