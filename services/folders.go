package services

import (
	"lain/data"
	"lain/jobs"
	"lain/repository"

	"github.com/gofiber/fiber/v2"
)

func GetFolders(userEmail, activeFolder string) []fiber.Map {
	count, _ := repository.CountFolders(userEmail)
	if count == 0 {
		jobs.SyncFolders(userEmail, data.FolderIcons)
	}

	return repository.BuildFolderTree(userEmail, activeFolder)
}

func GetFolderDisplayName(userEmail, folderPath string) string {
	return repository.GetFolderDisplayName(userEmail, folderPath)
}
