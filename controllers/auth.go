package controllers

import (
	"lain/config"
	"lain/repository"
	"lain/session"
	"lain/types"
	"lain/utils/crypto"
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

func Login(context *fiber.Ctx) error {
	var formData types.LoginForm
	if err := context.BodyParser(&formData); err != nil {
		return BadRequest(context, err)
	}

	encryptedPassword, err := crypto.Encrypt(formData.Password)
	if err != nil {
		return InternalServerError(context, err)
	}
	formData.Password = encryptedPassword

	preferences, err := repository.GetPreferences(formData)
	if err != nil {
		return InternalServerError(context, err)
	}

	if err = session.CreateSession(context, preferences.Email); err != nil {
		return InternalServerError(context, err)
	}

	return shortcuts.Redirect(context, "mail.inbox")
}

func Logout(context *fiber.Ctx) error {
	if err := session.DestroySession(context); err != nil {
		return InternalServerError(context, err)
	}

	return shortcuts.Redirect(context, "auth.login")
}
