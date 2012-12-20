package rrc

import . "go-cache"
import "go-cache/base"
import "time"

type RRCache struct {
	*base.BaseCache
	
	//cdb lists
	cdbl *CDBList

	Total int64
	Count int64
}

func NewRRCache(size int) *RRCache {
	c:= &RRCache{}
	c.BaseCache = base.NewBaseCache(size)
	c.cdbl = newCdbList(size)
	return c
}

func (c *RRCache) Get(key string) (object CacheObject, err error) {
	start := time.Now()
	tmp, err := c.get(key)
	entry := tmp.GetEntry()
	t := time.Since(start)
	c.Total += t.Nanoseconds()
	c.Count += 1
	if err == CacheMiss {
		var err1 error
		object, err1 = c.FetchFunc(key)
		entry.SetObject(object, c.CleanFunc)
		if err1 != nil {
			err = err1
		}
	} else {
		object = entry.GetObject()
	}
	return
}

func (c *RRCache) Set(key string, object CacheObject) {
	tmp, _ := c.get(key)
	entry := tmp.GetEntry()
	entry.SetObject(object, c.CleanFunc)
}

//get a CDB by a key
//in case of CacheMiss, the object stores in the cache entry is no longer valid
func (c *RRCache) get(key string) (base.CacheDirectoryBlock, error) {
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
	}
	if tmp.IsEntryNil() {
		panic("cannot be nil")
	}
	return tmp, err
}
