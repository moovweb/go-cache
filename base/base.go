package base

import . "go-cache"

func NewBaseCache(size int) *BaseCache {
	cache := &BaseCache{}
	cache.Size = size
	cache.CdbHash = make(map[string]CacheDirectoryBlock)
	cache.isGoroutineSafe = false
	cache.NewCacheEntryFunc = NewCacheEntry
	return cache
}

func NewSafeBaseCache(size int) *BaseCache {
	cache := &BaseCache{}
	cache.Size = size
	cache.CdbHash = make(map[string]CacheDirectoryBlock)
	cache.isGoroutineSafe = true
	cache.NewCacheEntryFunc = NewSafeCacheEntry
	return cache
}

func (c *BaseCache) SetFetchFunc(f CacheFetchFunc) {
	c.FetchFunc = f
}

func (c *BaseCache) SetCleanFunc(f CacheCleanFunc) {
	c.CleanFunc = f
}

func (c *BaseCache) GetAllObjects() map[string]CacheObject {
	all := make(map[string]CacheObject)
	for key, cdb := range(c.CdbHash) {
		if cdb.IsEntryCached() {
			entry := cdb.GetEntry()
			all[key] = entry.GetObject()
		}
	}
	return all
}

func (c *BaseCache) Fetch(key string) (CacheObject, error) {
	if c.FetchFunc == nil {
		return nil, CacheMiss
	}
	return c.FetchFunc(key)
}

func (c *BaseCache) GetHitRate() int {
	c.Lock()
	defer c.Unlock()
	if c.Accesses <= 0 {
		return 0
	}
	return int(c.Hits*100/c.Accesses)
}

func (c *BaseCache) Lock() {
	if c.isGoroutineSafe {
		c.mutex.Lock()
	}
}

func (c *BaseCache) Unlock() {
	if c.isGoroutineSafe {
		c.mutex.Unlock()
	}
}

func (c *BaseCache) CheckCache() {
	for key, cdb := range(c.CdbHash) {
		if cdb.GetKey() != key {
			panic("keys don't match")
		}
	}
}