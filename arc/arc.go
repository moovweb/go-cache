package arc

import . "go-cache"

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
}

func NewArcCache(size int) *ARCache {
	arc := &ARCache{}
	arc.size = size
	arc.t1 = newCdbList()
	arc.b1 = newCdbList()
	arc.t2 = newCdbList()
	arc.b2 = newCdbList()
	arc.cdbHash = make(map[string]*cacheDirectoryBlock)
	return arc
}

func (arc *ARCache) replace() *cacheEntry {
	var cdb *cacheDirectoryBlock
	if arc.t1.Len() >= max(1, arc.target_t1) {
		cdb = arc.t1.RemoveLRU()
		cdb.where = in_b1
		arc.b1.InsertMRU(cdb)
	} else {
		cdb = arc.t2.RemoveLRU()
		cdb.where = in_b2
		arc.b2.InsertMRU(cdb)
	}
	if cdb == nil || cdb.pointer == nil {
		panic("cdb is nil or cdb.pointer is nil")
	}
	p := cdb.pointer
	cdb.pointer = nil
	return p
}

func (arc *ARCache) Get(key string) (object CacheObject, err error) {
	tmp, err := arc.get(key)
	if err == CacheMiss {
		object, err = arc.fetchFunc(key)
		arc.setObject(tmp.pointer, object)
	} else {
		object = tmp.pointer.object
	}
	return
}

func (arc *ARCache) Set(key string, object CacheObject) {
	tmp, _ := arc.get(key)
	arc.setObject(tmp.pointer, object)
}

//get a CDB by a key
//in case of CacheMiss, the object stores in the cache entry is no longer valid
func (arc *ARCache) get(key string) (*cacheDirectoryBlock, error) {
	tmp := arc.cdbHash[key]
	var err error
	if tmp != nil {
		if tmp.where == in_t1 || tmp.where == in_t2 {
			if tmp.where == in_t1 {
				arc.t1.RemoveIt(tmp)
				arc.t2.InsertMRU(tmp)
				tmp.where = in_t2
			} else {
				arc.t2.SetMRU(tmp)
			}
		} else { //in b1 or b2
			if tmp.where == in_b1 {
				arc.target_t1 = min(arc.target_t1 + max(arc.b2.Len()/arc.b1.Len(), 1), arc.size)
				arc.b1.RemoveIt(tmp)
			} else {
				arc.target_t1 = min(arc.target_t1 - max(arc.b1.Len()/arc.b2.Len(), 1), 0)
				arc.b2.RemoveIt(tmp)
			}
			tmp.pointer = arc.replace()
			tmp.where = in_t2
			arc.t2.InsertMRU(tmp)
			err = CacheMiss
		}
	} else {
		if arc.t1.Len() + arc.b1.Len() == arc.size {
			if arc.t1.Len() < arc.size {
				tmp = arc.b1.RemoveLRU()
				tmp.pointer = arc.replace()
			} else {
				tmp = arc.t1.RemoveLRU()
			}
		} else if arc.t1.Len() + arc.t2.Len() + arc.b1.Len() + arc.b2.Len() >= arc.size {
			if arc.t1.Len() + arc.t2.Len() + arc.b1.Len() + arc.b2.Len() == arc.size * 2 {
				tmp = arc.b2.RemoveLRU()
			} else {
				tmp = newCacheDirectorBlock()
			}
			tmp.pointer = arc.replace()
		} else {
			tmp = newCacheDirectorBlock()
			tmp.pointer = newCacheEntry()
		}
		if len(tmp.key) > 0 {
			delete(arc.cdbHash, tmp.key)
		}
		tmp.key = key
		tmp.where = in_t1
		arc.t1.InsertMRU(tmp)
		arc.cdbHash[key] = tmp
		err = CacheMiss
	}
	if tmp.pointer == nil {
		panic("cannot be nil")
	}
	return tmp, err
}

func (arc *ARCache) SetFetchFunc(f CacheFetchFunc) {
	arc.fetchFunc = f
}

func (arc *ARCache) SetCleanFunc(f CacheCleanFunc) {
	arc.cleanFunc = f
}

func (arc *ARCache) GetAllObjects() map[string]CacheObject {
	all := make(map[string]CacheObject)
	for key, cdb := range(arc.cdbHash) {
		if cdb.where == in_t1 || cdb.where == in_t2 {
			all[key] = cdb.pointer.object
		}
	}
	return all
}

func (arc *ARCache) clearObject(entry *cacheEntry) {
	arc.setObject(entry, nil)
}

func (arc *ARCache) setObject(entry *cacheEntry, obj interface{}) {
	if entry != nil {
		if entry.object != nil && arc.cleanFunc != nil {
			arc.cleanFunc(entry.object)
		}
		entry.object = obj
	}
}

func (arc *ARCache) fetch(key string) (CacheObject, error) {
	if arc.fetchFunc == nil {
		return nil, CacheMiss
	}
	return arc.fetchFunc(key)
}

func (arc *ARCache) CheckCache() {
	for key, cdb := range(arc.cdbHash) {
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
