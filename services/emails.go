package services

import (
	"lain/jobs"
	"lain/models"
	"lain/repository"
	"lain/utils/crypto"
	"lain/utils/email"
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

func GetEmailDetails(userEmail string, emailID uint) (fiber.Map, error) {
	message, err := repository.GetEmailByID(userEmail, emailID)
	if err != nil {
		return nil, err
	}

	// Get attachments
	attachments, _ := repository.GetAttachmentsByEmailID(emailID)

	var attachmentMaps []fiber.Map
	for _, att := range attachments {
		attachmentMaps = append(attachmentMaps, fiber.Map{
			"ID":          att.ID,
			"Filename":    att.Filename,
			"ContentType": att.ContentType,
			"Size":        format.FormatFileSize(att.Size),
		})
	}

	// Sanitize HTML body
	body := message.BodyHTML
	if body == "" {
		body = "<pre>" + format.DecodeHTML(message.BodyText) + "</pre>"
	} else {
		body = format.SanitizeHTML(body)
	}

	return fiber.Map{
		"ID":          message.ID,
		"Subject":     format.DecodeHTML(message.Subject),
		"From":        format.DecodeHTML(message.From),
		"FromName":    format.DecodeHTML(message.FromName),
		"To":          format.DecodeHTML(message.To),
		"CC":          format.DecodeHTML(message.CC),
		"Date":        message.Date,
		"Body":        body,
		"IsRead":      message.IsRead,
		"IsFlagged":   message.IsFlagged,
		"Attachments": attachmentMaps,
	}, nil
}

func ToggleEmailFlag(userEmail string, emailID uint) (bool, error) {
	message, err := repository.GetEmailByID(userEmail, emailID)
	if err != nil {
		return false, err
	}

	prefs, err := repository.GetPreferencesByEmail(userEmail)
	if err != nil {
		return false, err
	}

	password, err := crypto.Decrypt(prefs.Authorization)
	if err != nil {
		return false, err
	}

	client, err := email.ConnectIMAP(userEmail, password)
	if err != nil {
		return false, err
	}
	defer email.DisconnectIMAP(client)

	if err := email.ToggleFlag(client, message.Folder.IMAPName, message.UID, message.IsFlagged); err != nil {
		return false, err
	}

	message.IsFlagged = !message.IsFlagged
	repository.UpdateEmail(message)

	return message.IsFlagged, nil
}

func MarkEmailAsRead(userEmail string, emailID uint) error {
	message, err := repository.GetEmailByID(userEmail, emailID)
	if err != nil {
		return err
	}

	if message.IsRead {
		return nil
	}

	prefs, err := repository.GetPreferencesByEmail(userEmail)
	if err != nil {
		return err
	}

	password, err := crypto.Decrypt(prefs.Authorization)
	if err != nil {
		return err
	}

	client, err := email.ConnectIMAP(userEmail, password)
	if err != nil {
		return err
	}
	defer email.DisconnectIMAP(client)

	if err := email.MarkAsRead(client, message.Folder.IMAPName, message.UID); err != nil {
		return err
	}

	message.IsRead = true
	return repository.UpdateEmail(message)
}
