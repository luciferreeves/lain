package router

import (
	"lain/controllers"
	"lain/types"
	"lain/utils/auth"
	"lain/utils/shortcuts"
	"lain/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func init() {
	urls.SetNamespace("auth")

	urls.Path(types.GET, "/login", func(c *fiber.Ctx) error {
		if auth.IsAuthenticated(c) {
			return shortcuts.Redirect(c, "mail.inbox")
		}
		return controllers.LoginPage(c)
	}, "login")

	urls.Path(types.GET, "/logout", controllers.Logout, "logout")
	urls.Path(types.POST, "/login", controllers.Login, "login.submit")
}
