package types

import (
	"time"

	"github.com/emersion/go-imap/client"
)

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

type EmailMessage struct {
	UID           uint32
	MessageID     string
	From          string
	FromName      string
	To            []string
	CC            []string
	BCC           []string
	ReplyTo       []string
	Subject       string
	Date          time.Time
	BodyText      string
	BodyHTML      string
	Size          uint32
	InReplyTo     string
	IsRead        bool
	IsFlagged     bool
	IsAnswered    bool
	IsDraft       bool
	HasAttachment bool
	Attachments   []EmailAttachment
}

type EmailAttachment struct {
	Filename    string
	ContentType string
	Data        []byte
}
