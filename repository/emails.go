package repository

import (
	"fmt"
	"lain/database"
	"lain/models"
)

func GetEmailsByFolder(userEmail string, folderID uint, limit int, offset int) ([]models.Email, error) {
	var messages []models.Email
	err := database.DB.Where("user_email = ? AND folder_id = ?", userEmail, folderID).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %w", err)
	}

	return messages, nil
}

func GetEmailByID(userEmail string, emailID uint) (*models.Email, error) {
	var message models.Email
	err := database.DB.Preload("Folder").Where("user_email = ? AND id = ?", userEmail, emailID).First(&message).Error
	if err != nil {
		return nil, fmt.Errorf("email not found: %w", err)
	}

	return &message, nil
}

func EmailExists(userEmail string, folderID uint, uid uint32) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Email{}).
		Where("user_email = ? AND folder_id = ? AND uid = ?", userEmail, folderID, uid).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateEmail(message *models.Email) (*models.Email, error) {
	if err := database.DB.Create(message).Error; err != nil {
		return nil, fmt.Errorf("failed to create email: %w", err)
	}
	return message, nil
}

func UpdateEmail(message *models.Email) error {
	if err := database.DB.Save(message).Error; err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}
	return nil
}

func GetAttachmentsByEmailID(emailID uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := database.DB.Where("email_id = ?", emailID).Find(&attachments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attachments: %w", err)
	}

	return attachments, nil
}

func CreateAttachment(attachment *models.Attachment) error {
	if err := database.DB.Create(attachment).Error; err != nil {
		return fmt.Errorf("failed to create attachment: %w", err)
	}
	return nil
}

func CountEmailsInFolder(userEmail string, folderID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Email{}).
		Where("user_email = ? AND folder_id = ?", userEmail, folderID).
		Count(&count).Error

	return count, err
}
