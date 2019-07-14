package schemautil

import (
	"fmt"
	"reflect"
)

func IsStruct(v reflect.Value) bool {
	return v.Type().Kind() == reflect.Struct
}

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

func Unpack(obj interface{}) []interface{} {
	out := make([]interface{}, 0)
	quene := []interface{}{obj}
	for len(quene) != 0 {
		var current interface{}
		current, quene = quene[0], quene[1:]
		out = append(out, current)
		t := reflect.ValueOf(current)
		if t.Kind() == reflect.Ptr {
			t = reflect.Indirect(t)
		}
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Kind() == reflect.Struct {
				quene = append(quene, t.Field(i).Interface())
			}
		}
	}
	return out
}

func floatToInt(val reflect.Value) int {
	return int(val.Float())
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return nil
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return fmt.Errorf("Provided value type didn't match obj field type (want %#v, got %#v)", structFieldType, val.Type())
	}
	structFieldValue.Set(val)
	return nil
}
