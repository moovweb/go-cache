package arc

import . "go-cache"
import "go-cache/base"

type ARCache struct {
	*base.BaseCache

	//cdb lists
	t1 *CDBList
	b1 *CDBList
	t2 *CDBList
	b2 *CDBList
	
	//param
	target_t1 int
}

func NewArcCache(size int) *ARCache {
	c := &ARCache{}
	c.BaseCache = base.NewBaseCache(size)
	c.t1 = newCdbList()
	c.b1 = newCdbList()
	c.t2 = newCdbList()
	c.b2 = newCdbList()
	return c
}


func NewSafeArcCache(size int) *ARCache {
	c := &ARCache{}
	c.BaseCache = base.NewSafeBaseCache(size)
	c.t1 = newCdbList()
	c.b1 = newCdbList()
	c.t2 = newCdbList()
	c.b2 = newCdbList()
	return c
}

func (c *ARCache) replace() *base.CacheEntry {
	var cdb base.CacheDirectoryBlock
	if c.t1.Len() >= max(1, c.target_t1) {
		cdb = c.t1.RemoveLRU()
		cdb.(*ArcCdb).where = in_b1
		c.b1.InsertMRU(cdb)
	} else {
		cdb = c.t2.RemoveLRU()
		cdb.(*ArcCdb).where = in_b2
		c.b2.InsertMRU(cdb)
	}
	if cdb == nil || cdb.IsEntryNil() {
		panic("cdb is nil or cdb.pointer is nil")
	}
	p := cdb.GetEntry()
	cdb.SetEntry(nil)
	return p
}

func (c *ARCache) Set(key string, object CacheObject) {
	tmp, _ := c.get(key)
	entry := tmp.GetEntry()
	entry.SetObject(object, c.CleanFunc)
}

//get a CDB by a key
//in case of CacheMiss, the object stores in the cache entry is no longer valid
func (c *ARCache) get(key string) (base.CacheDirectoryBlock, error) {
	c.Lock()
	defer c.Unlock()
	tmp := c.CdbHash[key]
	var err error
	if tmp != nil {
		if tmp.(*ArcCdb).where == in_t1 || tmp.(*ArcCdb).where == in_t2 {
			if tmp.(*ArcCdb).where == in_t1 {
				c.t1.RemoveIt(tmp)
				c.t2.InsertMRU(tmp)
				tmp.(*ArcCdb).where = in_t2
			} else {
				c.t2.SetMRU(tmp)
			}
			c.Hits += 1
		} else { //in b1 or b2
			if tmp.(*ArcCdb).where == in_b1 {
				c.target_t1 = min(c.target_t1 + max(c.b2.Len()/c.b1.Len(), 1), c.Size)
				c.b1.RemoveIt(tmp)
			} else {
				c.target_t1 = max(c.target_t1 - max(c.b1.Len()/c.b2.Len(), 1), 0)
				c.b2.RemoveIt(tmp)
			}
			tmp.SetEntry(c.replace())
			tmp.(*ArcCdb).where = in_t2
			c.t2.InsertMRU(tmp)
			err = CacheMiss
		}
	} else {
		if c.t1.Len() + c.b1.Len() == c.Size {
			if c.t1.Len() < c.Size {
				tmp = c.b1.RemoveLRU()
				tmp.SetEntry(c.replace())
			} else {
				tmp = c.t1.RemoveLRU()
			}
		} else if c.t1.Len() + c.t2.Len() + c.b1.Len() + c.b2.Len() >= c.Size {
			if c.t1.Len() + c.t2.Len() + c.b1.Len() + c.b2.Len() == c.Size * 2 {
				tmp = c.b2.RemoveLRU()
			} else {
				tmp = newCacheDirectorBlock()
			}
			tmp.SetEntry(c.replace())
		} else {
			tmp = newCacheDirectorBlock()
			tmp.SetEntry(c.NewCacheEntryFunc())
		}
		if len(tmp.GetKey()) > 0 {
			delete(c.CdbHash, tmp.GetKey())
		}
		tmp.SetKey(key)
		tmp.(*ArcCdb).where = in_t1
		c.t1.InsertMRU(tmp)
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

func (c *ARCache) CheckCache() {
	for key, cdb := range(c.CdbHash) {
		if cdb.GetKey() != key {
			panic("keys don't match")
		}
		if (cdb.(*ArcCdb).where == in_b1 || cdb.(*ArcCdb).where == in_b2) && !cdb.IsEntryNil() {
			panic("cdb pointer should be nil")
		}
		if (cdb.(*ArcCdb).where == in_t1 || cdb.(*ArcCdb).where == in_t2) && cdb.IsEntryNil() {
			panic("cdb pointer should not be nil")
		}
	}
}
