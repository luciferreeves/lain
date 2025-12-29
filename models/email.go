package models

import (
	"time"

	"gorm.io/gorm"
)

type Email struct {
	ID            uint   `gorm:"primaryKey"`
	UserEmail     string `gorm:"index:idx_user_folder,priority:1;not null"`
	FolderID      uint   `gorm:"index:idx_user_folder,priority:2;not null"`
	UID           uint32 `gorm:"not null"`
	MessageID     string `gorm:"index"`
	From          string `gorm:"not null"`
	FromName      string
	To            string `gorm:"not null"`
	CC            string
	BCC           string
	ReplyTo       string
	Subject       string
	Date          time.Time `gorm:"index"`
	BodyText      string    `gorm:"type:text"`
	BodyHTML      string    `gorm:"type:text"`
	RawHeaders    string    `gorm:"type:text"`
	Snippet       string
	Size          int64
	InReplyTo     string
	IsRead        bool `gorm:"default:false"`
	IsFlagged     bool `gorm:"default:false"`
	IsAnswered    bool `gorm:"default:false"`
	IsDraft       bool `gorm:"default:false"`
	HasAttachment bool `gorm:"default:false"`

	Folder      Folder       `gorm:"foreignKey:FolderID"`
	Attachments []Attachment `gorm:"foreignKey:EmailID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
