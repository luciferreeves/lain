package processors

import (
	"lain/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func request(ctx *fiber.Ctx) error {
	ctx.Locals("Request", meta.BuildRequest(ctx))
	return ctx.Next()
}
