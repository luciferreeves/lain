package shortcuts

import (
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, template string, data any) error {
	bind := make(fiber.Map)

	if err := mergeFlash(ctx, bind); err != nil {
		return err
	}

	mergeUserValues(ctx, bind)

	if data != nil {
		if err := mergeData(bind, data); err != nil {
			return err
		}
	}

	return ctx.Render(template, bind)
}

func RenderWithStatus(ctx *fiber.Ctx, template string, data any, statusCode int) error {
	ctx.Status(statusCode)
	return Render(ctx, template, data)
}
