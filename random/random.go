package rrc

import . "go-cache"
import "go-cache/base"

type RRCache struct {
	*base.BaseCache
	
	//cdb lists
	cdbl *CDBList
}

func NewRRCache(size int) *RRCache {
	c:= &RRCache{}
	c.BaseCache = base.NewBaseCache(size)
	c.cdbl = newCdbList(size)
	return c
}

func NewSafeRRCache(size int) *RRCache {
	c:= &RRCache{}
	c.BaseCache = base.NewSafeBaseCache(size)
	c.cdbl = newCdbList(size)
	return c
}

func (c *RRCache) Set(key string, object CacheObject) {
	tmp, _ := c.get(key)
	entry := tmp.GetEntry()
	entry.SetObject(object, c.CleanFunc)
}

//get a CDB by a key
//in case of CacheMiss, the object stores in the cache entry is no longer valid
func (c *RRCache) get(key string) (base.CacheDirectoryBlock, error) {
	c.Lock()
	defer c.Unlock()
	tmp := c.CdbHash[key]
	var err error
	if tmp == nil {
		if c.cdbl.Len() == c.Size {
			tmp = c.cdbl.Select()
		} else {
			tmp = newCacheDirectorBlock()
			tmp.SetEntry(base.NewCacheEntry())
			c.cdbl.Add(tmp)
		}
		if len(tmp.GetKey()) > 0 {
			delete(c.CdbHash, tmp.GetKey())
		}
		tmp.SetKey(key)
		c.CdbHash[key] = tmp
		err = CacheMiss
	} else {
		c.Hits += 1
	}
	c.Accesses += 1
	if tmp.IsEntryNil() {
		panic("cannot be nil")
	}
	return tmp, err
}
