package jobs

import (
	"lain/models"
	"lain/repository"
	"lain/types"
	"lain/utils/crypto"
	"lain/utils/email"
	"strings"
)

func SyncFolders(userEmail string, iconMap map[string]types.FolderIconVariant) error {
	prefs, err := repository.GetPreferencesByEmail(userEmail)
	if err != nil {
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
		if email.IsVirtualFolder(imapFolder.Name) {
			continue
		}

		existingFolder, err := repository.GetFolderByIMAPName(userEmail, imapFolder.Name)

		sortOrder := email.GetSortOrder(imapFolder.Name, i)
		folderType := email.GetFolderType(imapFolder.Name, iconMap)
		iconVariant := iconMap[folderType]

		if err != nil {
			folder := models.Folder{
				UserEmail: userEmail,
				Name:      imapFolder.Name,
				IMAPName:  imapFolder.Name,
				IconOpen:  iconVariant.Open,
				IconClose: iconVariant.Close,
				SortOrder: sortOrder,
			}
			created, err := repository.CreateFolder(&folder)
			if err != nil {
				continue
			}
			foldersByName[strings.ToLower(imapFolder.Name)] = created.ID
		} else {
			existingFolder.Name = imapFolder.Name
			existingFolder.SortOrder = sortOrder
			existingFolder.IconOpen = iconVariant.Open
			existingFolder.IconClose = iconVariant.Close
			repository.UpdateFolder(existingFolder)
			foldersByName[strings.ToLower(imapFolder.Name)] = existingFolder.ID
		}
	}

	for _, imapFolder := range imapFolders {
		if email.IsVirtualFolder(imapFolder.Name) {
			continue
		}

		if strings.Contains(imapFolder.Name, "/") {
			parts := strings.Split(imapFolder.Name, "/")
			if len(parts) > 1 {
				parentName := strings.Join(parts[:len(parts)-1], "/")
				parentNameLower := strings.ToLower(parentName)
				if parentID, ok := foldersByName[parentNameLower]; ok {
					folder, err := repository.GetFolderByIMAPName(userEmail, imapFolder.Name)
					if err == nil {
						folder.ParentID = &parentID
						repository.UpdateFolder(folder)
					}
				}
			}
		}
	}

	return nil
}
