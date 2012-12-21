// +build !cache_debug

package base

import . "go-cache"
import "sync"

type BaseCache struct {
	//max number of cache entries
	Size int

	//the hook which is called in case of cache miss
	FetchFunc CacheFetchFunc

	//clean func
	CleanFunc CacheCleanFunc

	//hash table for searching cache entries
	CdbHash map[string]CacheDirectoryBlock

	Accesses int64
	Hits int64

	isGoroutineSafe bool
	mutex sync.Mutex

	NewCacheEntryFunc func()*CacheEntry 
}

func (c *BaseCache) PrintStats() {
	println("cache hit rate:", c.GetHitRate())
}