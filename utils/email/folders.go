package email

import (
	"lain/types"

	"github.com/emersion/go-imap"
)

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
