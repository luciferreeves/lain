package services

import (
	"lain/jobs"
	"lain/models"
	"lain/repository"
	"lain/utils/format"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetEmails(userEmail, folderPath string, prefs *models.Preferences, page int) ([]fiber.Map, error) {
	decodedPath, _ := url.QueryUnescape(folderPath)

	folder, err := repository.GetFolderByIMAPName(userEmail, strings.ToLower(decodedPath))
	if err != nil {
		return nil, err
	}

	emailCount, _ := repository.CountEmailsInFolder(userEmail, folder.ID)
	if emailCount == 0 {
		jobs.SyncEmails(userEmail, folder.ID, folder.IMAPName)
	}

	limit := prefs.EmailsPerPage
	offset := (page - 1) * limit

	messages, err := repository.GetEmailsByFolder(userEmail, folder.ID, limit, offset)
	if err != nil {
		return []fiber.Map{}, err
	}

	var emails []fiber.Map
	for _, message := range messages {
		fromName := message.FromName
		if fromName == "" {
			fromName = message.From
		}

		emails = append(emails, fiber.Map{
			"ID":            message.ID,
			"UID":           message.UID,
			"From":          format.DecodeHTML(message.From),
			"FromName":      format.DecodeHTML(fromName),
			"Subject":       format.DecodeHTML(message.Subject),
			"Date":          message.Date,
			"DateFormatted": format.FormatEmailDate(message.Date, prefs.DateFormat, prefs.TimeFormat, prefs.PrettyDates, prefs.TimeZone),
			"Snippet":       format.DecodeHTML(message.Snippet),
			"IsRead":        message.IsRead,
			"IsFlagged":     message.IsFlagged,
			"HasAttachment": message.HasAttachment,
		})
	}

	return emails, nil
}
