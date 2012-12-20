package base

import . "go-cache"

type BaseCache struct {
	//max number of cache entries
	Size int

	//the hook which is called in case of cache miss
	FetchFunc CacheFetchFunc

	//clean func
	CleanFunc CacheCleanFunc

	//hash table for searching cache entries
	CdbHash map[string]CacheDirectoryBlock
}

func NewBaseCache(size int) *BaseCache {
	cache := &BaseCache{}
	cache.Size = size
	cache.CdbHash = make(map[string]CacheDirectoryBlock)
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

func (c *BaseCache) CheckCache() {
	for key, cdb := range(c.CdbHash) {
		if cdb.GetKey() != key {
			panic("keys don't match")
		}
	}
}