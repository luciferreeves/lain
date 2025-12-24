package jobs

import (
	"fmt"
	"lain/database"
	"lain/models"
	"lain/repository"
	"lain/utils/crypto"
	"lain/utils/email"
	"lain/utils/format"
	"lain/utils/storage"
	"strings"
)

func SyncEmails(userEmail string, folderID uint, folderPath string) error {
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
		exists, err := repository.EmailExists(userEmail, folderID, msg.UID)
		if err != nil {
			continue
		}
		if exists {
			continue
		}

		snippet := format.GenerateSnippet(msg.BodyText, msg.BodyHTML)

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

		createdMessage, err := repository.CreateEmail(&message)
		if err != nil {
			continue
		}

		for _, att := range msg.Attachments {
			path, err := storage.UploadAttachment(userEmail, createdMessage.ID, att.Filename, att.Data, att.ContentType)
			if err != nil {
				continue
			}

			attachment := models.Attachment{
				EmailID:     createdMessage.ID,
				Filename:    att.Filename,
				ContentType: att.ContentType,
				Size:        int64(len(att.Data)),
				MinIOPath:   path,
			}

			repository.CreateAttachment(&attachment)
		}
	}

	return nil
}
