package processors

import (
	"lain/repository"
	"lain/session"
	"lain/utils/auth"

	"github.com/gofiber/fiber/v2"
)

func preferences(ctx *fiber.Ctx) error {
	if auth.IsAuthenticated(ctx) {
		email, err := session.GetSessionEmail(ctx)
		if err == nil {
			prefs, err := repository.GetPreferencesByEmail(email)
			if err == nil {
				ctx.Locals("Preferences", prefs)
			}
		}
	}

	return ctx.Next()
}
