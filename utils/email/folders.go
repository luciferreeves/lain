package email

import (
	"lain/models"
	"lain/types"
	"net/url"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/gofiber/fiber/v2"
)

func GetDisplayName(imapName string) string {
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

func FetchFolders(client *types.EmailClient) ([]types.IMAPFolder, error) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- client.List("", "*", mailboxes)
	}()

	var folders []types.IMAPFolder
	for m := range mailboxes {
		folders = append(folders, types.IMAPFolder{
			Name:        m.Name,
			HasChildren: hasAttribute(m.Attributes, "\\HasChildren"),
		})
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return folders, nil
}

func GetFolderType(folderName string, iconMap map[string]types.FolderIconVariant) string {
	nameLower := strings.ToLower(folderName)

	if strings.Contains(folderName, "/") {
		parts := strings.Split(folderName, "/")
		nameLower = strings.ToLower(parts[len(parts)-1])
	}

	for iconType := range iconMap {
		if iconType == "default" {
			continue
		}
		if strings.Contains(nameLower, iconType) {
			return iconType
		}
	}

	return "default"
}

func SortFolders(folders []models.Folder) {
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

func GetSortOrder(folderName string, index int) int {
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
		baseOrder := GetSortOrder(parts[0], index)
		return baseOrder + 1000 + (index * 10)
	}

	return 100 + index
}

func CopyFolderMap(folder fiber.Map) fiber.Map {
	copy := fiber.Map{}
	for k, v := range folder {
		if k == "Subfolders" {
			if subfolders, ok := v.([]fiber.Map); ok {
				subfoldersCopy := make([]fiber.Map, len(subfolders))
				for i, sf := range subfolders {
					subfoldersCopy[i] = CopyFolderMap(sf)
				}
				copy[k] = subfoldersCopy
			}
		} else {
			copy[k] = v
		}
	}
	return copy
}

func IsVirtualFolder(folderName string) bool {
	return strings.HasPrefix(folderName, "Virtual") || strings.Contains(folderName, "/Virtual")
}

func UpdateActiveFolder(folders []fiber.Map, activeFolder string) []fiber.Map {
	decodedActive, _ := url.QueryUnescape(activeFolder)
	activeLower := strings.ToLower(decodedActive)

	foldersCopy := make([]fiber.Map, len(folders))
	for i, f := range folders {
		foldersCopy[i] = CopyFolderMap(f)
	}

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

	updateActive(foldersCopy)
	return foldersCopy
}
