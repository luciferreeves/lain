package repository

import (
	"fmt"
	"lain/cache"
	"lain/database"
	"lain/models"
	"lain/types"
	"lain/utils/email"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var FolderIcons = map[string]types.FolderIconVariant{
	"default": {
		Open:  "/static/icons/folder_open.png",
		Close: "/static/icons/folder.png",
	},
	"inbox": {
		Open:  "/static/icons/inbox_open.png",
		Close: "/static/icons/inbox.png",
	},
	"encrypt": {
		Open:  "/static/icons/encrypt_open.png",
		Close: "/static/icons/encrypt.png",
	},
	"dog": {
		Open:  "/static/icons/dog_open.png",
		Close: "/static/icons/dog.png",
	},
	"internal": {
		Open:  "/static/icons/internal_open.png",
		Close: "/static/icons/internal.png",
	},
	"draft": {
		Open:  "/static/icons/draft_open.png",
		Close: "/static/icons/draft.png",
	},
	"progress": {
		Open:  "/static/icons/draft_open.png",
		Close: "/static/icons/draft.png",
	},
	"sent": {
		Open:  "/static/icons/sent.png",
		Close: "/static/icons/sent.png",
	},
	"archive": {
		Open:  "/static/icons/archive_open.png",
		Close: "/static/icons/archive.png",
	},
	"trash": {
		Open:  "/static/icons/trash_open.png",
		Close: "/static/icons/trash.png",
	},
	"delete": {
		Open:  "/static/icons/trash_open.png",
		Close: "/static/icons/trash.png",
	},
	"spam": {
		Open:  "/static/icons/junk_open.png",
		Close: "/static/icons/junk.png",
	},
	"junk": {
		Open:  "/static/icons/junk_open.png",
		Close: "/static/icons/junk.png",
	},
}

func GetAllFolders(userEmail string) ([]models.Folder, error) {
	var folders []models.Folder
	err := database.DB.Where("user_email = ?", userEmail).Find(&folders).Error
	return folders, err
}

func GetFolderByIMAPName(userEmail, imapName string) (*models.Folder, error) {
	var folder models.Folder
	err := database.DB.Where("user_email = ? AND LOWER(imap_name) = ?", userEmail, strings.ToLower(imapName)).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func CreateFolder(folder *models.Folder) (*models.Folder, error) {
	if err := database.DB.Create(folder).Error; err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}
	return folder, nil
}

func UpdateFolder(folder *models.Folder) error {
	if err := database.DB.Save(folder).Error; err != nil {
		return fmt.Errorf("failed to update folder: %w", err)
	}
	return nil
}

func CountFolders(userEmail string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Folder{}).Where("user_email = ?", userEmail).Count(&count).Error
	return count, err
}

func BuildFolderTree(userEmail, activeFolder string) []fiber.Map {
	if cached, ok := cache.GetFolders(userEmail); ok {
		return email.UpdateActiveFolder(cached, activeFolder)
	}

	allFolders, _ := GetAllFolders(userEmail)
	email.SortFolders(allFolders)

	folderMap := make(map[uint]*fiber.Map)
	var rootFolders []fiber.Map

	for _, folder := range allFolders {
		displayName := email.GetDisplayName(folder.IMAPName)

		folderData := fiber.Map{
			"ID":          folder.ID,
			"Name":        displayName,
			"IMAPName":    folder.IMAPName,
			"IconOpen":    folder.IconOpen,
			"IconClose":   folder.IconClose,
			"UnreadCount": folder.UnreadCount,
			"Active":      false,
			"ParentID":    folder.ParentID,
			"SortOrder":   folder.SortOrder,
			"Subfolders":  []fiber.Map{},
		}
		folderMap[folder.ID] = &folderData
	}

	for _, folder := range allFolders {
		folderData := folderMap[folder.ID]
		if folder.ParentID == nil {
			rootFolders = append(rootFolders, *folderData)
		} else {
			if parent, ok := folderMap[*folder.ParentID]; ok {
				subfolders := (*parent)["Subfolders"].([]fiber.Map)
				(*parent)["Subfolders"] = append(subfolders, *folderData)
			}
		}
	}

	cache.SetFolders(userEmail, rootFolders)
	return email.UpdateActiveFolder(rootFolders, activeFolder)
}

func GetFolderDisplayName(userEmail, folderPath string) string {
	decodedPath, _ := url.QueryUnescape(folderPath)

	folder, err := GetFolderByIMAPName(userEmail, decodedPath)
	if err != nil {
		if strings.ToLower(decodedPath) == "inbox" {
			return "Inbox"
		}
		return decodedPath
	}

	return email.GetDisplayName(folder.IMAPName)
}
