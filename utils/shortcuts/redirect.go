package shortcuts

import (
	"lain/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func Redirect(ctx *fiber.Ctx, routeName string) error {
	path, ok := urls.GetFullPath(routeName)
	if !ok {
		return fiber.ErrNotFound
	}
	return ctx.Redirect(path)
}

func RedirectWithStatus(ctx *fiber.Ctx, routeName string, statusCode int) error {
	path, ok := urls.GetFullPath(routeName)
	if !ok {
		return fiber.ErrNotFound
	}
	return ctx.Redirect(path, statusCode)
}
