package router

import (
	"lain/types"
	"lain/utils/auth"
	"lain/utils/shortcuts"
	"lain/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func init() {
	urls.SetNamespace("")

	urls.Path(types.GET, "/", func(c *fiber.Ctx) error {
		if auth.IsAuthenticated(c) {
			return shortcuts.Redirect(c, "mail.inbox")
		}
		return shortcuts.Redirect(c, "auth.login")
	}, "home")
}
