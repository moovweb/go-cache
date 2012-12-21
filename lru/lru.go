package lru

import . "go-cache"
import "go-cache/base"

type LRUCache struct {
	*base.BaseCache
	
	//cdb lists
	cdbl *CDBList
}

func NewLRUCache(size int) *LRUCache {
	c := &LRUCache{}
	c.BaseCache = base.NewBaseCache(size)
	c.cdbl = newCdbList()
	return c
}

func NewSafeLRUCache(size int) *LRUCache {
	c := &LRUCache{}
	c.BaseCache = base.NewSafeBaseCache(size)
	c.cdbl = newCdbList()
	return c
}

func (c *LRUCache) Set(key string, object CacheObject) {
	tmp, _ := c.get(key)
	entry := tmp.GetEntry()
	entry.SetObject(object, c.CleanFunc)
}

//get a CDB by a key
//in case of CacheMiss, the object stores in the cache entry is no longer valid
func (c *LRUCache) get(key string) (base.CacheDirectoryBlock, error) {
	c.Lock()
	defer c.Unlock()
	tmp := c.CdbHash[key]
	var err error
	if tmp != nil {
		c.cdbl.SetMRU(tmp)
		c.Hits += 1
	} else {
		if c.cdbl.Len() == c.Size {
			tmp = c.cdbl.RemoveLRU()
		} else {
			tmp = newCacheDirectorBlock()
			tmp.SetEntry(c.NewCacheEntryFunc())
		}
		if len(tmp.GetKey()) > 0 {
			delete(c.CdbHash, tmp.GetKey())
		}
		tmp.SetKey(key)
		c.cdbl.InsertMRU(tmp)
		c.CdbHash[key] = tmp
		err = CacheMiss
	}
	c.Accesses += 1
	if tmp.IsEntryNil() {
		panic("cannot be nil")
	}
	if err == CacheMiss {
		tmp.GetEntry().SetDirty()
	}
	return tmp, err
}