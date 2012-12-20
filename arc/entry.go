package arc

import . "go-cache"
import "sync"

type cacheEntry struct {
	object CacheObject
	rwlock  sync.RWMutex
}

func newCacheEntry() *cacheEntry {
	return &cacheEntry{}
}