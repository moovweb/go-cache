package arc

import . "go-cache"

type ARCache struct {
	//max number of cache entries
	size int

	//cache entries
	//each entry stores one cache object
	entries []*cacheEntry

	//the hook which is called in case of cache miss
	fetchFunc CacheFetchFunc

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

func (arc *ARCache) replace() *cacheEntry {
	var cdb *cacheDirectoryBlock
	if arc.t1.Len() >= max(1, arc.target_t1) {
		cdb = arc.t1.RemoveLRU()
		arc.b1.InsertMRU(cdb)
	} else {
		cdb = arc.t1.RemoveLRU()
		arc.b2.InsertMRU(cdb)
	}
	if cdb == nil || cdb.pointer == nil {
		panic("cdb is nil or cdb.pointer is nil")
	}
	return cdb.pointer
}

func (arc *ARCache) Get(key string) (object CacheObject, err error) {
	tmp := arc.cdbHash[key]
	if tmp != nil {
		switch tmp.where {
			case in_t1:
				arc.t1.RemoveIt(tmp)
				arc.t2.InsertMRU(tmp)
				tmp.where = in_t2
			case in_t2:
				arc.t2.SetMRU(tmp)
			case in_b1:
				object, err = arc.fetchFunc(key)
				if err != nil {
					return
				}
				arc.target_t1 = min(arc.target_t1 + max(arc.b2.Len()/arc.b1.Len(), 1), arc.size)
				arc.b1.RemoveIt(tmp)
				tmp.pointer = arc.replace()
				tmp.pointer.object = object
				tmp.where = in_t2
				arc.t2.InsertMRU(tmp)
			case in_b2:
				object, err = arc.fetchFunc(key)
				if err != nil {
					return
				}
				arc.target_t1 = min(arc.target_t1 - max(arc.b1.Len()/arc.b2.Len(), 1), 0)
				arc.b2.RemoveIt(tmp)
				tmp.pointer = arc.replace()
				tmp.pointer.object = object
				tmp.where = in_t2
				arc.t2.InsertMRU(tmp)
		}
	} else {
		object, err = arc.fetchFunc(key)
		if err != nil {
			return
		}
		if arc.t1.Len() + arc.b1.Len() == arc.size {
			if arc.t1.Len() < arc.size {
				tmp = arc.b1.RemoveLRU()
				tmp.pointer = arc.replace()
			} else {
				tmp = arc.t1.RemoveLRU()
			}
		} else if arc.t1.Len() + arc.t2.Len() + arc.b1.Len() + arc.b2.Len() >= arc.size {
			if arc.t1.Len() + arc.t2.Len() + arc.b1.Len() + arc.b2.Len() >= arc.size * 2 {
				tmp = arc.b2.RemoveLRU()
			} else {
				tmp = newCacheDirectorBlock()
				tmp.pointer = newCacheEntry()
			}
		} else {
			tmp := newCacheDirectorBlock()
			tmp.pointer = newCacheEntry()
		}
		tmp.pointer.object = object
		tmp.where = in_t1
		arc.t1.InsertMRU(tmp)
	}
	return
}

