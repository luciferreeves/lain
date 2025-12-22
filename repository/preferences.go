package repository

import (
	"lain/database"
	"lain/models"
	"lain/types"

	"gorm.io/gorm"
)

func GetPreferences(formData types.LoginForm) (*models.Preferences, error) {
	var preferences models.Preferences

	if err := database.DB.Where("email = ?", formData.Email).First(&preferences).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return CreateDefaultPreferences(formData)
		}
		return nil, err
	}

	return &preferences, nil
}

func CreateDefaultPreferences(formData types.LoginForm) (*models.Preferences, error) {
	preferences := models.Preferences{
		Email:         formData.Email,
		Authorization: formData.Password,
	}

	if err := database.DB.Create(&preferences).Error; err != nil {
		return nil, err
	}

	return &preferences, nil
}
