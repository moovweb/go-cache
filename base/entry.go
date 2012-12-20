package base

import . "go-cache"

type CacheEntry struct {
	object CacheObject
}

func NewCacheEntry() *CacheEntry {
	return &CacheEntry{}
}