package repository

import (
	"fmt"
	"lain/database"
	"lain/models"
	"lain/utils/crypto"
	"lain/utils/email"
	"lain/utils/storage"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetEmails(userEmail, folderPath string, limit int, offset int) ([]fiber.Map, error) {
	// Decode URL-encoded path (e.g., "inbox/lets%20encrypt" -> "inbox/lets encrypt")
	decodedPath, _ := url.QueryUnescape(folderPath)

	var folder models.Folder
	err := database.DB.Where("user_email = ? AND LOWER(imap_name) = ?", userEmail, strings.ToLower(decodedPath)).First(&folder).Error
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}

	var count int64
	database.DB.Model(&models.Email{}).Where("user_email = ? AND folder_id = ?", userEmail, folder.ID).Count(&count)

	// Always sync if no emails exist
	if count == 0 {
		if err := syncEmails(userEmail, folder.ID, folder.IMAPName); err != nil {
			// Log the error but continue to show UI
			fmt.Printf("Failed to sync emails for folder %s: %v\n", folder.IMAPName, err)
		}
		// Recount after sync
		database.DB.Model(&models.Email{}).Where("user_email = ? AND folder_id = ?", userEmail, folder.ID).Count(&count)
	}

	var messages []models.Email
	err = database.DB.Where("user_email = ? AND folder_id = ?", userEmail, folder.ID).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %w", err)
	}

	var emailMaps []fiber.Map
	for _, message := range messages {
		emailMaps = append(emailMaps, fiber.Map{
			"ID":            message.ID,
			"UID":           message.UID,
			"From":          message.From,
			"FromName":      message.FromName,
			"Subject":       message.Subject,
			"Date":          message.Date,
			"Snippet":       message.Snippet,
			"IsRead":        message.IsRead,
			"IsFlagged":     message.IsFlagged,
			"HasAttachment": message.HasAttachment,
		})
	}

	return emailMaps, nil
}

func GetEmail(userEmail string, emailID uint) (*models.Email, error) {
	var message models.Email
	err := database.DB.Preload("Folder").Where("user_email = ? AND id = ?", userEmail, emailID).First(&message).Error
	if err != nil {
		return nil, fmt.Errorf("email not found: %w", err)
	}

	return &message, nil
}

func GetAttachments(emailID uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := database.DB.Where("email_id = ?", emailID).Find(&attachments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attachments: %w", err)
	}

	return attachments, nil
}

func MarkEmailAsRead(userEmail string, emailID uint) error {
	var message models.Email
	err := database.DB.Preload("Folder").Where("user_email = ? AND id = ?", userEmail, emailID).First(&message).Error
	if err != nil {
		return fmt.Errorf("email not found: %w", err)
	}

	if message.IsRead {
		return nil
	}

	var prefs models.Preferences
	if err := database.DB.Where("email = ?", userEmail).First(&prefs).Error; err != nil {
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
	if err := database.DB.Save(&message).Error; err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}

func ToggleEmailFlag(userEmail string, emailID uint) error {
	var message models.Email
	err := database.DB.Preload("Folder").Where("user_email = ? AND id = ?", userEmail, emailID).First(&message).Error
	if err != nil {
		return fmt.Errorf("email not found: %w", err)
	}

	var prefs models.Preferences
	if err := database.DB.Where("email = ?", userEmail).First(&prefs).Error; err != nil {
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

	if err := email.ToggleFlag(client, message.Folder.IMAPName, message.UID, message.IsFlagged); err != nil {
		return err
	}

	message.IsFlagged = !message.IsFlagged
	if err := database.DB.Save(&message).Error; err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}

func syncEmails(userEmail string, folderID uint, folderPath string) error {
	var prefs models.Preferences
	if err := database.DB.Where("email = ?", userEmail).First(&prefs).Error; err != nil {
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

	messages, err := email.FetchMessages(client, folderPath, 50)
	if err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	for _, msg := range messages {
		var existingMessage models.Email
		result := database.DB.Where("user_email = ? AND folder_id = ? AND uid = ?", userEmail, folderID, msg.UID).First(&existingMessage)

		if result.Error == nil {
			continue
		}

		snippet := generateSnippet(msg.BodyText, msg.BodyHTML)

		message := models.Email{
			UserEmail:     userEmail,
			FolderID:      folderID,
			UID:           msg.UID,
			MessageID:     msg.MessageID,
			From:          msg.From,
			FromName:      msg.FromName,
			To:            strings.Join(msg.To, ", "),
			CC:            strings.Join(msg.CC, ", "),
			BCC:           strings.Join(msg.BCC, ", "),
			ReplyTo:       strings.Join(msg.ReplyTo, ", "),
			Subject:       msg.Subject,
			Date:          msg.Date,
			BodyText:      msg.BodyText,
			BodyHTML:      msg.BodyHTML,
			Snippet:       snippet,
			Size:          int64(msg.Size),
			InReplyTo:     msg.InReplyTo,
			IsRead:        msg.IsRead,
			IsFlagged:     msg.IsFlagged,
			IsAnswered:    msg.IsAnswered,
			IsDraft:       msg.IsDraft,
			HasAttachment: msg.HasAttachment,
		}

		if err := database.DB.Create(&message).Error; err != nil {
			continue
		}

		for _, att := range msg.Attachments {
			path, err := storage.UploadAttachment(userEmail, message.ID, att.Filename, att.Data, att.ContentType)
			if err != nil {
				continue
			}

			attachment := models.Attachment{
				EmailID:     message.ID,
				Filename:    att.Filename,
				ContentType: att.ContentType,
				Size:        int64(len(att.Data)),
				MinIOPath:   path,
			}

			database.DB.Create(&attachment)
		}
	}

	return nil
}

func generateSnippet(bodyText, bodyHTML string) string {
	text := bodyText
	if text == "" && bodyHTML != "" {
		text = stripHTML(bodyHTML)
	}

	text = strings.TrimSpace(text)
	if len(text) > 150 {
		text = text[:150] + "..."
	}

	return text
}

func stripHTML(html string) string {
	text := html
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n\n")
	text = strings.ReplaceAll(text, "</div>", "\n")

	inTag := false
	var result strings.Builder
	for _, char := range text {
		if char == '<' {
			inTag = true
			continue
		}
		if char == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(char)
		}
	}

	return strings.TrimSpace(result.String())
}
