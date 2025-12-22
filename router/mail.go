package router

import (
	"lain/session"
	"lain/types"
	"lain/utils/auth"
	"lain/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func init() {
	urls.SetNamespace("mail")

	urls.Path(types.GET, "/inbox", auth.RequireAuthentication(func(c *fiber.Ctx) error {
		email, _ := session.GetSessionEmail(c)
		return c.SendString("Inbox for " + email)
	}), "inbox")
}
