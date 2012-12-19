package lru

import . "go-cache"

type LRUCache struct {
	//max number of cache entries
	size int

	//cache entries
	//each entry stores one cache object
	//entries []*cacheEntry

	//the hook which is called in case of cache miss
	fetchFunc CacheFetchFunc

	//clean func
	cleanFunc CacheCleanFunc

	//cdb lists
	cdbl *CDBList
	
	//cdb hash table for searching cdb
	cdbHash map[string]*cacheDirectoryBlock
}

func NewLRUCache(size int) *LRUCache {
	cache := &LRUCache{}
	cache.size = size
	cache.cdbl = newCdbList()
	cache.cdbHash = make(map[string]*cacheDirectoryBlock)
	return cache
}

func (c *LRUCache) Get(key string) (object CacheObject, err error) {
	tmp, err := c.get(key)
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

func (c *LRUCache) Set(key string, object CacheObject) {
	tmp, _ := c.get(key)
	c.setObject(tmp.pointer, object)
}

//get a CDB by a key
//in case of CacheMiss, the object stores in the cache entry is no longer valid
func (c *LRUCache) get(key string) (*cacheDirectoryBlock, error) {
	tmp := c.cdbHash[key]
	var err error
	if tmp != nil {
		c.cdbl.SetMRU(tmp)
	} else {
		if c.cdbl.Len() == c.size {
			tmp = c.cdbl.RemoveLRU()
		} else {
			tmp = newCacheDirectorBlock()
			tmp.pointer = newCacheEntry()
		}
		if len(tmp.key) > 0 {
			delete(c.cdbHash, tmp.key)
		}
		tmp.key = key
		c.cdbl.InsertMRU(tmp)
		c.cdbHash[key] = tmp
		err = CacheMiss
	}
	if tmp.pointer == nil {
		panic("cannot be nil")
	}
	return tmp, err
}

func (c *LRUCache) SetFetchFunc(f CacheFetchFunc) {
	c.fetchFunc = f
}

func (c *LRUCache) SetCleanFunc(f CacheCleanFunc) {
	c.cleanFunc = f
}

func (c *LRUCache) GetAllObjects() map[string]CacheObject {
	all := make(map[string]CacheObject)
	for key, cdb := range(c.cdbHash) {
		all[key] = cdb.pointer.object
	}
	return all
}

func (c *LRUCache) clearObject(entry *cacheEntry) {
	c.setObject(entry, nil)
}

func (c *LRUCache) setObject(entry *cacheEntry, obj interface{}) {
	if entry != nil {
		if entry.object != nil && c.cleanFunc != nil {
			c.cleanFunc(entry.object)
		}
		entry.object = obj
	}
}

func (c *LRUCache) fetch(key string) (CacheObject, error) {
	if c.fetchFunc == nil {
		return nil, CacheMiss
	}
	return c.fetchFunc(key)
}

func (c *LRUCache) CheckCache() {
	for key, cdb := range(c.cdbHash) {
		if cdb.key != key {
			panic("keys don't match")
		}
	}
}