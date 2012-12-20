package arc

import . "go-cache"
import "time"
import "sync"

type ARCache struct {
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
	t1 *CDBList
	b1 *CDBList
	t2 *CDBList
	b2 *CDBList
	
	//cdb hash table for searching cdb
	cdbHash map[string]*cacheDirectoryBlock

	//param
	target_t1 int

	//mutex
	mutex sync.Mutex

	Total int64
	Count int64
}

func NewArcCache(size int) *ARCache {
	c := &ARCache{}
	c.size = size
	c.t1 = newCdbList()
	c.b1 = newCdbList()
	c.t2 = newCdbList()
	c.b2 = newCdbList()
	c.cdbHash = make(map[string]*cacheDirectoryBlock)
	return c
}

func (c *ARCache) replace() *cacheEntry {
	var cdb *cacheDirectoryBlock
	if c.t1.Len() >= max(1, c.target_t1) {
		cdb = c.t1.RemoveLRU()
		cdb.where = in_b1
		c.b1.InsertMRU(cdb)
	} else {
		cdb = c.t2.RemoveLRU()
		cdb.where = in_b2
		c.b2.InsertMRU(cdb)
	}
	if cdb == nil || cdb.pointer == nil {
		panic("cdb is nil or cdb.pointer is nil")
	}
	p := cdb.pointer
	cdb.pointer = nil
	return p
}

func (c *ARCache) Get(key string) (object CacheObject, err error) {
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

func (c *ARCache) Set(key string, object CacheObject) {
	tmp, _ := c.get(key)
	c.setObject(tmp.pointer, object)
}

//get a CDB by a key
//in case of CacheMiss, the object stores in the cache entry is no longer valid
func (c *ARCache) get(key string) (*cacheDirectoryBlock, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp := c.cdbHash[key]
	var err error
	if tmp != nil {
		if tmp.where == in_t1 || tmp.where == in_t2 {
			if tmp.where == in_t1 {
				c.t1.RemoveIt(tmp)
				c.t2.InsertMRU(tmp)
				tmp.where = in_t2
			} else {
				c.t2.SetMRU(tmp)
			}
		} else { //in b1 or b2
			if tmp.where == in_b1 {
				c.target_t1 = min(c.target_t1 + max(c.b2.Len()/c.b1.Len(), 1), c.size)
				c.b1.RemoveIt(tmp)
			} else {
				c.target_t1 = max(c.target_t1 - max(c.b1.Len()/c.b2.Len(), 1), 0)
				c.b2.RemoveIt(tmp)
			}
			tmp.pointer = c.replace()
			tmp.where = in_t2
			c.t2.InsertMRU(tmp)
			err = CacheMiss
		}
	} else {
		if c.t1.Len() + c.b1.Len() == c.size {
			if c.t1.Len() < c.size {
				tmp = c.b1.RemoveLRU()
				tmp.pointer = c.replace()
			} else {
				tmp = c.t1.RemoveLRU()
			}
		} else if c.t1.Len() + c.t2.Len() + c.b1.Len() + c.b2.Len() >= c.size {
			if c.t1.Len() + c.t2.Len() + c.b1.Len() + c.b2.Len() == c.size * 2 {
				tmp = c.b2.RemoveLRU()
			} else {
				tmp = newCacheDirectorBlock()
			}
			tmp.pointer = c.replace()
		} else {
			tmp = newCacheDirectorBlock()
			tmp.pointer = newCacheEntry()
		}
		if len(tmp.key) > 0 {
			delete(c.cdbHash, tmp.key)
		}
		tmp.key = key
		tmp.where = in_t1
		c.t1.InsertMRU(tmp)
		c.cdbHash[key] = tmp
		err = CacheMiss
	}
	if tmp.pointer == nil {
		panic("cannot be nil")
	}
	return tmp, err
}

func (c *ARCache) SetFetchFunc(f CacheFetchFunc) {
	c.fetchFunc = f
}

func (c *ARCache) SetCleanFunc(f CacheCleanFunc) {
	c.cleanFunc = f
}

func (c *ARCache) GetAllObjects() map[string]CacheObject {
	all := make(map[string]CacheObject)
	for key, cdb := range(c.cdbHash) {
		if cdb.where == in_t1 || cdb.where == in_t2 {
			all[key] = cdb.pointer.object
		}
	}
	return all
}

func (c *ARCache) clearObject(entry *cacheEntry) {
	c.setObject(entry, nil)
}

func (c *ARCache) setObject(entry *cacheEntry, obj CacheObject) {
	if entry != nil {
		if entry.object != nil && c.cleanFunc != nil {
			c.cleanFunc(entry.object)
		}
		entry.object = obj
	}
}

func (c *ARCache) fetch(key string) (CacheObject, error) {
	if c.fetchFunc == nil {
		return nil, CacheMiss
	}
	return c.fetchFunc(key)
}

func (c *ARCache) CheckCache() {
	for key, cdb := range(c.cdbHash) {
		if cdb.key != key {
			panic("keys don't match")
		}
		if (cdb.where == in_b1 || cdb.where == in_b2) && cdb.pointer != nil {
			panic("cdb pointer should be nil")
		}
		if (cdb.where == in_t1 || cdb.where == in_t2) && cdb.pointer == nil {
			panic("cdb pointer should not be nil")
		}
	}
}
