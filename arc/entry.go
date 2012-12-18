package arc

import . "go-cache"

type cacheEntry struct {
	object CacheObject
}

var count = 0

func newCacheEntry() *cacheEntry {
	count += 1
	return &cacheEntry{}
}