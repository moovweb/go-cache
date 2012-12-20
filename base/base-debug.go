// +build cache_debug

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

	AccessTime int64
}

func (c *BaseCache) AddAccessTime(t int64) {
	c.AccessTime += t
}

func (c *BaseCache) GetAvgAccessTime() int64 {
	c.Lock()
	defer c.Unlock()
	if c.Accesses == 0 {
		return 0
	}
	return c.AccessTime/c.Accesses
}