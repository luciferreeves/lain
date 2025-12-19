package middleware

import (
	"lain/utils/auth"
	"lain/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func authentication(context *fiber.Ctx) error {
	if !auth.IsAuthenticated(context) {
		return shortcuts.Redirect(context, "auth.login")
	}
	return context.Next()
}
