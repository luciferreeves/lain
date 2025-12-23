package models

import (
	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model
	UserEmail string `gorm:"index"`

	Name      string
	IMAPName  string
	IconOpen  string
	IconClose string
	ParentID  *uint

	UnreadCount int `gorm:"default:0"`
	TotalCount  int `gorm:"default:0"`
	SortOrder   int `gorm:"default:0"`
}
