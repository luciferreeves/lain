package controllers

import (
	"lain/config"
	"lain/utils/meta"
	"lain/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func LoginPage(context *fiber.Ctx) error {
	meta.SetPageTitle(context, "Login")

	return shortcuts.Render(context, "auth/login", fiber.Map{
		"AllowedDomains": config.Server.AllowedDomains,
	})
}
