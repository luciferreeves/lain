package processors

import (
	"lain/config"

	"github.com/gofiber/fiber/v2"
)

const defaultTitle = "Lain | Present day, present time!"

func metadata(ctx *fiber.Ctx) error {
	ctx.Locals("Title", defaultTitle)
	ctx.Locals("AppName", config.Server.AppName)
	ctx.Locals("AppDescription", config.Server.AppDescription)
	ctx.Locals("AppEngine", config.Server.AppEngine)

	return ctx.Next()
}
