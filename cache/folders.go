package cache

import (
	"lain/types"
	"time"

	"github.com/gofiber/fiber/v2"
)

var folders *types.FolderCache

func init() {
	folders = &types.FolderCache{
		Data: make(map[string]*types.FolderCacheEntry),
		TTL:  5 * time.Minute,
	}
}

func GetFolders(userEmail string) ([]fiber.Map, bool) {
	folders.Mu.RLock()
	defer folders.Mu.RUnlock()

	entry, exists := folders.Data[userEmail]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Folders, true
}

func SetFolders(userEmail string, folderList []fiber.Map) {
	folders.Mu.Lock()
	defer folders.Mu.Unlock()

	now := time.Now()
	folders.Data[userEmail] = &types.FolderCacheEntry{
		Folders:   folderList,
		CachedAt:  now,
		ExpiresAt: now.Add(folders.TTL),
	}
}

func InvalidateFolders(userEmail string) {
	folders.Mu.Lock()
	defer folders.Mu.Unlock()

	delete(folders.Data, userEmail)
}

func InvalidateAllFolders() {
	folders.Mu.Lock()
	defer folders.Mu.Unlock()

	folders.Data = make(map[string]*types.FolderCacheEntry)
}

func SetFolderTTL(duration time.Duration) {
	folders.Mu.Lock()
	defer folders.Mu.Unlock()

	folders.TTL = duration
}
