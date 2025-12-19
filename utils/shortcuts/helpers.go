package shortcuts

import (
	"fmt"
	"reflect"
	"strings"
)

func structValue(data any) (reflect.Value, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf(
			"Render: unsupported bind type %T; must be struct or *struct",
			data,
		)
	}

	return v, nil
}

func mapStruct(v reflect.Value) map[string]any {
	t := v.Type()
	result := make(map[string]any, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		key := fieldType.Name
		if tag := fieldType.Tag.Get("json"); tag != "" && tag != "-" {
			if idx := strings.IndexByte(tag, ','); idx >= 0 {
				if idx > 0 {
					key = tag[:idx]
				}
			} else {
				key = tag
			}
		}

		result[key] = v.Field(i).Interface()
	}

	return result
}
