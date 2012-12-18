package arc

type cacheEntry struct {
	object CacheObject
}

func newCacheEntry() *cacheEntry {
	return &cacheEntry{}
}