package controllers

import (
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

	folders := repository.GetFolders(email, folderPath)
	displayName := repository.GetFolderDisplayName(email, folderPath)

	emails := []fiber.Map{}

	meta.SetPageTitle(context, displayName)

	return shortcuts.Render(context, "mail/folder", fiber.Map{
		"Folders": folders,
		"Emails":  emails,
		"Email":   nil,
	})
}
