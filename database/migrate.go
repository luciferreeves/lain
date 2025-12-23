package database

import "lain/models"

func migrate() error {
	err := DB.AutoMigrate(
		&models.Preferences{},
		&models.Folder{},
		&models.Email{},
	)
	return err
}
