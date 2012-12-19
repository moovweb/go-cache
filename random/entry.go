package rrc

import . "go-cache"

type cacheEntry struct {
	object CacheObject
	cdb *cacheDirectoryBlock
}

func newCacheEntry() *cacheEntry {
	return &cacheEntry{}
}