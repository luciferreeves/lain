package controllers

import (
	"lain/models"
	"lain/services"
	"lain/session"
	"lain/utils/meta"
	"lain/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Mailbox(context *fiber.Ctx) error {
	folderPath := context.Params("*", "inbox")
	if folderPath == "" {
		folderPath = "inbox"
	}

	userEmail, err := session.GetSessionEmail(context)
	if err != nil {
		return InternalServerError(context, err)
	}

	prefs := context.Locals("Preferences").(*models.Preferences)

	folders := services.GetFolders(userEmail, folderPath)
	displayName := services.GetFolderDisplayName(userEmail, folderPath)

	page := context.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	emails, err := services.GetEmails(userEmail, folderPath, prefs, page)
	if err != nil {
		emails = []fiber.Map{}
	}

	meta.SetPageTitle(context, displayName)

	return shortcuts.Render(context, "mail/folder", fiber.Map{
		"Folders": folders,
		"Emails":  emails,
		"Email":   nil,
		"Page":    page,
	})
}
