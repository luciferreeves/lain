package types

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type FolderCacheEntry struct {
	Folders   []fiber.Map
	CachedAt  time.Time
	ExpiresAt time.Time
}

type FolderCache struct {
	Mu   sync.RWMutex
	Data map[string]*FolderCacheEntry
	TTL  time.Duration
}
