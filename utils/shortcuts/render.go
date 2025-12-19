package shortcuts

import (
	"maps"

	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, template string, data any) error {
	bind := make(fiber.Map)

	ctx.Context().VisitUserValues(func(key []byte, value any) {
		bind[string(key)] = value
	})

	if data != nil {
		switch v := data.(type) {
		case map[string]any:
			maps.Copy(bind, v)
		case fiber.Map:
			maps.Copy(bind, v)
		default:
			rv, err := structValue(data)
			if err != nil {
				return err
			}

			maps.Copy(bind, mapStruct(rv))
		}
	}

	return ctx.Render(template, bind)
}

func RenderWithStatus(ctx *fiber.Ctx, template string, data any, statusCode int) error {
	ctx.Status(statusCode)
	return Render(ctx, template, data)
}
