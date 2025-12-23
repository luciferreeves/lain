package repository

import (
	"lain/database"
	"lain/models"
	"lain/types"
	"lain/utils/crypto"
	"lain/utils/email"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var folderIcons = map[string]types.FolderIconVariant{
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

func GetFolders(userEmail, activeFolder string) []fiber.Map {
	syncFolders(userEmail)

	var allFolders []models.Folder
	database.DB.Where("user_email = ?", userEmail).Find(&allFolders)

	sortFolders(allFolders)

	folderMap := make(map[uint]*fiber.Map)
	var rootFolders []fiber.Map

	for _, folder := range allFolders {
		displayName := getDisplayName(folder.IMAPName)

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

	return updateActiveFolder(rootFolders, activeFolder)
}

func GetFolderDisplayName(userEmail, folderPath string) string {
	decodedPath, _ := url.QueryUnescape(folderPath)

	var folder models.Folder
	err := database.DB.Where("user_email = ? AND LOWER(imap_name) = ?", userEmail, strings.ToLower(decodedPath)).First(&folder).Error

	if err != nil {
		if strings.ToLower(decodedPath) == "inbox" {
			return "Inbox"
		}
		return decodedPath
	}

	return getDisplayName(folder.IMAPName)
}

func getDisplayName(imapName string) string {
	if strings.Contains(imapName, "/") {
		parts := strings.Split(imapName, "/")
		lastPart := parts[len(parts)-1]
		if strings.ToLower(lastPart) == "inbox" {
			return "Inbox"
		}
		return lastPart
	}

	if strings.ToLower(imapName) == "inbox" {
		return "Inbox"
	}

	return imapName
}

func updateActiveFolder(folders []fiber.Map, activeFolder string) []fiber.Map {
	decodedActive, _ := url.QueryUnescape(activeFolder)
	activeLower := strings.ToLower(decodedActive)

	var updateActive func([]fiber.Map)
	updateActive = func(folderList []fiber.Map) {
		for i := range folderList {
			imapNameLower := strings.ToLower(folderList[i]["IMAPName"].(string))
			folderList[i]["Active"] = imapNameLower == activeLower
			if subfolders, ok := folderList[i]["Subfolders"].([]fiber.Map); ok && len(subfolders) > 0 {
				updateActive(subfolders)
			}
		}
	}

	updateActive(folders)
	return folders
}

func sortFolders(folders []models.Folder) {
	for i := 0; i < len(folders)-1; i++ {
		for j := 0; j < len(folders)-i-1; j++ {
			if folders[j].SortOrder > folders[j+1].SortOrder {
				folders[j], folders[j+1] = folders[j+1], folders[j]
			} else if folders[j].SortOrder == folders[j+1].SortOrder {
				if strings.ToLower(folders[j].IMAPName) > strings.ToLower(folders[j+1].IMAPName) {
					folders[j], folders[j+1] = folders[j+1], folders[j]
				}
			}
		}
	}
}

func getFolderType(folderName string) string {
	nameLower := strings.ToLower(folderName)

	if strings.Contains(folderName, "/") {
		parts := strings.Split(folderName, "/")
		nameLower = strings.ToLower(parts[len(parts)-1])
	}

	for iconType := range folderIcons {
		if iconType == "default" {
			continue
		}
		if strings.Contains(nameLower, iconType) {
			return iconType
		}
	}

	return "default"
}

func syncFolders(userEmail string) error {
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

	imapFolders, err := email.FetchFolders(client)
	if err != nil {
		return err
	}

	foldersByName := make(map[string]uint)

	for i, imapFolder := range imapFolders {
		var folder models.Folder
		imapNameLower := strings.ToLower(imapFolder.Name)
		result := database.DB.Where("user_email = ? AND LOWER(imap_name) = ?", userEmail, imapNameLower).First(&folder)

		sortOrder := getSortOrder(imapFolder.Name, i)
		folderType := getFolderType(imapFolder.Name)
		iconVariant := folderIcons[folderType]

		if result.Error != nil {
			folder = models.Folder{
				UserEmail: userEmail,
				Name:      imapFolder.Name,
				IMAPName:  imapFolder.Name,
				IconOpen:  iconVariant.Open,
				IconClose: iconVariant.Close,
				SortOrder: sortOrder,
			}
			database.DB.Create(&folder)
			foldersByName[imapNameLower] = folder.ID
		} else {
			folder.Name = imapFolder.Name
			folder.SortOrder = sortOrder
			folder.IconOpen = iconVariant.Open
			folder.IconClose = iconVariant.Close
			database.DB.Save(&folder)
			foldersByName[imapNameLower] = folder.ID
		}
	}

	for _, imapFolder := range imapFolders {
		if strings.Contains(imapFolder.Name, "/") {
			parts := strings.Split(imapFolder.Name, "/")
			if len(parts) > 1 {
				parentName := strings.Join(parts[:len(parts)-1], "/")
				parentNameLower := strings.ToLower(parentName)
				if parentID, ok := foldersByName[parentNameLower]; ok {
					var folder models.Folder
					imapNameLower := strings.ToLower(imapFolder.Name)
					if err := database.DB.Where("user_email = ? AND LOWER(imap_name) = ?", userEmail, imapNameLower).First(&folder).Error; err == nil {
						folder.ParentID = &parentID
						database.DB.Save(&folder)
					}
				}
			}
		}
	}

	return nil
}

func getSortOrder(folderName string, index int) int {
	nameLower := strings.ToLower(folderName)

	if nameLower == "inbox" {
		return 0
	}
	if strings.Contains(nameLower, "draft") {
		return 1
	}
	if strings.Contains(nameLower, "sent") {
		return 2
	}
	if strings.Contains(nameLower, "archive") {
		return 3
	}
	if strings.Contains(nameLower, "trash") || strings.Contains(nameLower, "deleted") {
		return 4
	}
	if strings.Contains(nameLower, "spam") || strings.Contains(nameLower, "junk") {
		return 5
	}

	if strings.Contains(folderName, "/") {
		parts := strings.Split(folderName, "/")
		baseOrder := getSortOrder(parts[0], index)
		return baseOrder + 1000 + (index * 10)
	}

	return 100 + index
}
