package controllers

import (
	"lain/models"
	"lain/repository"
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

	email, err := session.GetSessionEmail(context)
	if err != nil {
		return InternalServerError(context, err)
	}

	prefs := context.Locals("Preferences").(*models.Preferences)

	folders := repository.GetFolders(email, folderPath)
	displayName := repository.GetFolderDisplayName(email, folderPath)

	page := context.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := prefs.EmailsPerPage
	offset := (page - 1) * limit

	emails, err := repository.GetEmails(email, folderPath, limit, offset)
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
