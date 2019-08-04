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

func InterfaceToSlice(i interface{}) ([]interface{}, bool) {
	if !IsIterable(i) {
		return nil, false
	}
	s := reflect.ValueOf(i)
	out := make([]interface{}, s.Len())
	for j := 0; j < s.Len(); j++ {
		out[j] = s.Index(j).Interface()
	}
	return out, true
}

func FlattenMap(m map[string]interface{}, delimter string) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := FlattenMap(child, delimter)
			for nk, nv := range nm {
				key := fmt.Sprintf("%s%s%s", k, delimter, nk)
				o[key] = nv
			}
		default:
			o[k] = v
		}
	}
	return o
}

func FlattenValueSlice(i interface{}) string {
	if v, ok := InterfaceToSlice(i); ok {
		out := make([]string, len(v))
		for j, val := range v {
			out[j] = fmt.Sprintf("%v", val)
		}
		return strings.Join(out, ",")
	}
	return fmt.Sprintf("%v", i)
}

func flattenKey(i interface{}, delimiter string) map[string]string {
	out := make(map[string]string)
	input := structs.Map(i)
	m := FlattenMap(input, delimiter)
	for k, v := range m {
		key := strings.ToLower(k)
		out[key] = FlattenValueSlice(v)
	}
	return out
}

func FlattenKey(i interface{}, delimiter string) []map[string]string {
	if v, ok := InterfaceToSlice(i); ok {
		out := make([]map[string]string, len(v))
		for j, val := range v {
			out[j] = flattenKey(val, delimiter)
		}
		return out
	}
	return []map[string]string{flattenKey(i, delimiter)}
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
