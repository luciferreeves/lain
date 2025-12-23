package models

import (
	"time"

	"gorm.io/gorm"
)

type Email struct {
	gorm.Model
	UserEmail string
	FolderID  uint
	Folder    Folder `gorm:"foreignKey:FolderID"`

	UID       uint32 `gorm:"index"`
	MessageID string `gorm:"index"`

	From     string
	FromName string
	To       string
	CC       string
	BCC      string
	ReplyTo  string

	Subject string
	Date    time.Time `gorm:"index"`

	BodyText string `gorm:"type:text"`
	BodyHTML string `gorm:"type:text"`
	Snippet  string

	IsRead        bool `gorm:"default:false;index"`
	IsFlagged     bool `gorm:"default:false;index"`
	IsAnswered    bool `gorm:"default:false"`
	IsDraft       bool `gorm:"default:false"`
	HasAttachment bool `gorm:"default:false;index"`

	Size       int64
	InReplyTo  string
	References string
}

type Attachment struct {
	gorm.Model
	EmailID uint
	Email   Email `gorm:"foreignKey:EmailID"`

	Filename    string
	ContentType string
	Size        int64
	ContentID   string
	MinIOPath   string `gorm:"index"`
}
