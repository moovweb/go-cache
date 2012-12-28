package base

import . "go-cache"
import "sync"

type BaseCache struct {
	//the limit of the total size of cached objects
	size int

	//this is called on the object evicted from the cache
	CleanFunc CacheCleanFunc

	//stats info
	//total number of accesses
	accesses int64
	//total number of hits
	hits int64

	//is this cache safe for multi-goroutines
	isGoroutineSafe bool
	//the mutex to make it goroutine safe
	mutex sync.Mutex

	//CacheDirectoryBlock Manager
	CdbManager
}

func NewBaseCache(size int, cdbm CdbManager) *BaseCache {
	cache := &BaseCache{}
	cache.size = size
	cache.isGoroutineSafe = false
	cache.CdbManager = cdbm
	return cache
}

func NewSafeBaseCache(size int, cdbm CdbManager) *BaseCache {
	cache := &BaseCache{}
	cache.size = size
	cache.isGoroutineSafe = true
	cache.CdbManager = cdbm
	return cache
}

func (c *BaseCache) SetCleanFunc(f CacheCleanFunc) {
	c.CleanFunc = f
}

func (c *BaseCache) Get(key string) (object CacheObject, err error) {
	if len(key) == 0 {
		return nil, EmptyKey
	}
	c.Lock()
	defer c.Unlock()
	c.accesses += 1
	cdb, err := c.CdbManager.Find(key)
	if err != nil {
		return nil, err
	}
	c.hits += 1
	return cdb.GetObject(), nil
}


func (c *BaseCache) Set(key string, object CacheObject) error {
	if len(key) == 0 {
		return EmptyKey
	}
	c.Lock()
	defer c.Unlock()
	return c.CdbManager.Replace(key, object, c.size, c.CleanFunc)
}

func (c *BaseCache) GetHitRate() int {
	c.Lock()
	defer c.Unlock()
	if c.accesses <= 0 {
		return 0
	}
	return int(c.hits*100/c.accesses)
}

func (c *BaseCache) GetUsage() int {
	c.Lock()
	defer c.Unlock()
	return c.CdbManager.GetUsage()
}

func (c *BaseCache) Check() {
	c.Lock()
	defer c.Unlock()
	c.CdbManager.Check()
}

func (c *BaseCache) Collect() map[string]CacheObject {
	c.Lock()
	defer c.Unlock()
	return c.CdbManager.Collect()
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