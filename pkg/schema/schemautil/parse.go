package schemautil

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/oliveagle/jsonpath"
)

const DefaultTagName = "jsonpath"

func parseJsonpath(in interface{}, out interface{}) (map[string]interface{}, error) {
	tags := GetTaggedFields(out, DefaultTagName)
	obj := make(map[string]interface{})
	for _, tag := range tags {
		tokens := strings.Split(tag.Tag, ",")
		pattern, _ := tokens[0], tokens[1:]
		jpath, err := jsonpath.Compile(pattern)
		if err != nil {
			return obj, err
		}
		value, err := jpath.Lookup(in)
		if err != nil {
			return obj, err
		}
		if IsStruct(tag.Value) {
			nested, err := parseJsonpath(value, tag.Value.Interface())
			if err != nil {
				return obj, err
			}
			obj[tag.Name] = nested
		} else {
			obj[tag.Name] = value
		}
	}
	return obj, nil
}

func ParseJsonpath(in interface{}, out interface{}) error {
	obj, err := parseJsonpath(in, out)
	if err != nil {
		return err
	}
	return mapstructure.WeakDecode(obj, out)
}
