package types

import "github.com/emersion/go-imap/client"

type EmailClient struct {
	*client.Client
}

type IMAPFolder struct {
	Name        string
	HasChildren bool
}

type FolderIconVariant struct {
	Open  string
	Close string
}
