package auth

import (
	session "lain/session"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(context *fiber.Ctx) bool {
	session, err := session.Store.Get(context)
	if err != nil {
		return false
	}

	email := session.Get("email")
	return email != nil
}
