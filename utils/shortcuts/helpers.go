package shortcuts

import (
	"fmt"
	"maps"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func structValue(data any) (reflect.Value, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("unsupported bind type %T; must be struct or *struct", data)
	}

	return v, nil
}

func mapStruct(v reflect.Value) map[string]any {
	t := v.Type()
	result := make(map[string]any, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		result[getFieldKey(field)] = v.Field(i).Interface()
	}

	return result
}

func getFieldKey(field reflect.StructField) string {
	key := field.Name
	tag := field.Tag.Get("json")

	if tag == "" || tag == "-" {
		return key
	}

	if idx := strings.IndexByte(tag, ','); idx >= 0 {
		if idx > 0 {
			return tag[:idx]
		}
		return key
	}

	return tag
}

func normalizeBind(data any) (fiber.Map, error) {
	if data == nil {
		return nil, nil
	}

	switch v := data.(type) {
	case fiber.Map:
		return v, nil
	case map[string]any:
		return fiber.Map(v), nil
	default:
		return structToMap(v)
	}
}

func structToMap(data any) (fiber.Map, error) {
	v, err := structValue(data)
	if err != nil {
		return nil, err
	}
	return fiber.Map(mapStruct(v)), nil
}

func mergeFlash(ctx *fiber.Ctx, bind fiber.Map) error {
	flash, err := ConsumeFlash(ctx)
	if err != nil || flash == nil {
		return err
	}

	flashMap, err := normalizeBind(flash)
	if err != nil {
		return err
	}

	maps.Copy(bind, flashMap)
	return nil
}

func mergeUserValues(ctx *fiber.Ctx, bind fiber.Map) {
	ctx.Context().VisitUserValues(func(key []byte, value any) {
		bind[string(key)] = value
	})
}

func mergeData(bind fiber.Map, data any) error {
	dataMap, err := normalizeBind(data)
	if err != nil {
		return err
	}

	maps.Copy(bind, dataMap)
	return nil
}
