package fifo

import . "go-cache"
import "time"

type BaseCache struct {
	//max number of cache entries
	Size int

	//cache entries
	//each entry stores one cache object
	//entries []*cacheEntry

	//the hook which is called in case of cache miss
	FetchFunc CacheFetchFunc

	//clean func
	CleanFunc CacheCleanFunc

	//cdb lists
	Cdbl *CdbList
	
	//cdb hash table for searching cdb
	CdbHash map[string]*cacheDirectoryBlock
}

func NewRRCache(size int) *RRCache {
	cache := &RRCache{}
	cache.size = size
	cache.cdbl = newCdbList(size)
	cache.cdbHash = make(map[string]*cacheDirectoryBlock)
	return cache
}

func (c *RRCache) Get(key string) (object CacheObject, err error) {
	start := time.Now()
	tmp, err := c.get(key)
	t := time.Since(start)
	c.Total += t.Nanoseconds()
	c.Count += 1
	if err == CacheMiss {
		var err1 error
		object, err1 = c.fetchFunc(key)
		c.setObject(tmp.pointer, object)
		if err1 != nil {
			err = err1
		}
	} else {
		object = tmp.pointer.object
	}
	return
}

func (c *RRCache) Set(key string, object CacheObject) {
	tmp, _ := c.get(key)
	c.setObject(tmp.pointer, object)
}

//get a CDB by a key
//in case of CacheMiss, the object stores in the cache entry is no longer valid
func (c *RRCache) get(key string) (*cacheDirectoryBlock, error) {
	tmp := c.cdbHash[key]
	var err error
	if tmp == nil {
		if c.cdbl.Len() == c.size {
			tmp = c.cdbl.RandomSelection()
		} else {
			tmp = newCacheDirectorBlock()
			tmp.pointer = newCacheEntry()
			c.cdbl.Add(tmp)
		}
		if len(tmp.key) > 0 {
			delete(c.cdbHash, tmp.key)
		}
		tmp.key = key
		c.cdbHash[key] = tmp
		err = CacheMiss
	}
	if tmp.pointer == nil {
		panic("cannot be nil")
	}
	return tmp, err
}

func (c *RRCache) SetFetchFunc(f CacheFetchFunc) {
	c.fetchFunc = f
}

func (c *RRCache) SetCleanFunc(f CacheCleanFunc) {
	c.cleanFunc = f
}

func (c *RRCache) GetAllObjects() map[string]CacheObject {
	all := make(map[string]CacheObject)
	for key, cdb := range(c.cdbHash) {
		all[key] = cdb.pointer.object
	}
	return all
}

func (c *RRCache) clearObject(entry *cacheEntry) {
	c.setObject(entry, nil)
}

func (c *RRCache) setObject(entry *cacheEntry, obj CacheObject) {
	if entry != nil {
		if entry.object != nil && c.cleanFunc != nil {
			c.cleanFunc(entry.object)
		}
		entry.object = obj
	}
}

func (c *RRCache) fetch(key string) (CacheObject, error) {
	if c.fetchFunc == nil {
		return nil, CacheMiss
	}
	return c.fetchFunc(key)
}

func (c *RRCache) CheckCache() {
	for key, cdb := range(c.cdbHash) {
		if cdb.key != key {
			panic("keys don't match")
		}
	}
}