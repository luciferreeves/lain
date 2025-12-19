package meta

import (
	"fmt"
	"lain/config"

	"github.com/gofiber/fiber/v2"
)

func SetPageTitle(context *fiber.Ctx, title string) {
	title = fmt.Sprintf("%s | %s", title, config.Server.AppName)
	context.Locals("Title", title)
}
