package database

import "lain/models"

func migrate() error {
	err := DB.AutoMigrate(
		&models.Preferences{},
	)
	return err
}
