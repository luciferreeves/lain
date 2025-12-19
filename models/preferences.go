package models

import (
	"lain/types"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Preferences struct {
	// Credentials
	Email         string `gorm:"primaryKey;uniqueIndex"`
	Authorization string `gorm:"not null"`

	// Calendar
	CalendarName                string
	CalendarURL                 string
	CalendarUsername            string
	CalendarAuthorization       string
	CalendarColor               string
	CalendarDescription         string
	CalendarSyncIntervalMinutes int `gorm:"default:30"`

	// Address Book
	AddressBookName                string
	AddressBookURL                 string
	AddressBookUsername            string
	AddressBookAuthorization       string
	AddressBookColor               string
	AddressBookDescription         string
	AddressBookSyncIntervalMinutes int `gorm:"default:30"`

	// Mailbox View
	Language           string                   `gorm:"default:'en'"`
	TimeZone           string                   `gorm:"default:'(GMT +00:00) UTC'"`
	TimeFormat         types.TimeFormat         `gorm:"type:varchar(20);default:'07:30'"`
	DateFormat         types.DateFormat         `gorm:"type:varchar(20);default:'12/20/2025'"`
	PrettyDates        bool                     `gorm:"default:true"`
	MarkMessagesAsRead types.EmailReadingOption `gorm:"type:varchar(20);default:'Immediately'"`
	EmailsPerPage      int                      `gorm:"default:50"`
	EnableSounds       bool                     `gorm:"default:false"`

	// Displaying Emails
	OpenMessagesInNewWindow             bool                               `gorm:"default:false"`
	ShowEmailAddressWithDisplayName     bool                               `gorm:"default:false"`
	DisplayHTML                         bool                               `gorm:"default:true"`
	LoadRemoteContent                   types.RemoteResourceDownloadOption `gorm:"type:varchar(30);default:'From my contacts'"`
	ReturnReceipts                      types.ReturnReceiptOption          `gorm:"type:varchar(50);default:'Ask me each time'"`
	DisplayAttachedImagesBelowMessage   bool                               `gorm:"default:true"`
	DisplayEmoticonsInPlainTextMessages bool                               `gorm:"default:false"`

	// Composing Emails
	ComposeMessagesInNewWindow   bool                                  `gorm:"default:false"`
	DefaultComposeFormat         types.EmailComposingOption            `gorm:"type:varchar(100);default:'Always, except when replying to plain text messages'"`
	AutoSaveDraftIntervalSeconds types.AutoSaveDraftIntervalOption     `gorm:"default:300"`
	AlwaysRequestReturnReceipt   bool                                  `gorm:"default:false"`
	AlwaysRequestDeliveryStatus  bool                                  `gorm:"default:false"`
	QuoteOriginalMessage         types.EmailReplyOption                `gorm:"type:varchar(50);default:'Place my reply below the original message'"`
	MessageForwarding            types.MessageForwardingOption         `gorm:"type:varchar(50);default:'Inline'"`
	HTMLFontFamily               types.EmailHTMLFontFamilyOption       `gorm:"type:varchar(50);default:'Verdana'"`
	HTMLFontSize                 types.EmailHTMLFontSizeOption         `gorm:"default:'10'"`
	EnableEmoticons              bool                                  `gorm:"default:true"`
	AttachmentNames              types.EmailAttachementNameStyleOption `gorm:"type:varchar(30);default:'RFC 2047/2231 (Outlook)'"`

	// Signature Options
	IncludeSignature                    types.EmailSignatureOption `gorm:"type:varchar(20);default:'Always'"`
	PlaceSignatureBelowQuotedText       bool                       `gorm:"default:false"`
	RemoveOriginalSignatureWhenReplying bool                       `gorm:"default:true"`

	// Contact Options
	ShowContactPhotos bool                       `gorm:"default:true"`
	ListContactsAs    types.ContactDisplayOption `gorm:"type:varchar(20);default:'Display Name'"`
	SortContactsBy    types.ContactSortingOption `gorm:"type:varchar(20);default:'Last Name'"`
	ContactsPerPage   int                        `gorm:"default:50"`

	// Server Settings
	MarkMessagesAsReadOnDelete             bool                           `gorm:"default:true"`
	FlagMessagesForDeletionInsteadOfDelete bool                           `gorm:"default:false"`
	DoNotShowDeletedMessages               bool                           `gorm:"default:false"`
	DirectlyDeleteMessagesInJunkFolder     bool                           `gorm:"default:false"`
	MarkMessagesAsReadOnArchive            bool                           `gorm:"default:true"`
	ClearTrashOnLogout                     types.ClearTrashOnLogoutOption `gorm:"type:varchar(30);default:'Never'"`

	UseEmailEncryption bool           `gorm:"default:false"`
	PGPPrivateKey      datatypes.JSON `gorm:"type:jsonb"`
	PGPPublicKey       datatypes.JSON `gorm:"type:jsonb"`

	// Sync
	EmailSyncIntervalMinutes int `gorm:"default:5"`
	LastSyncedAt             time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
