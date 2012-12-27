package arc

import (
	. "go-cache"
	"go-cache/base"
	"container/list"
)

const (
	in_t1 = iota
	in_b1
	in_t2
	in_b2
	not_in
)

type ArcCdbm struct {
	*base.BasicCdbm
	t1 *CdbList
	b1 *CdbList
	t2 *CdbList
	b2 *CdbList
	target_t1 int
}


type ArcCdb struct {
	element *list.Element
	where int
	base.BaseCdb
	size int
}

func newCacheDirectorBlock() *ArcCdb {
	cdb := &ArcCdb{}
	return cdb
}

func newCdbm() *ArcCdbm {
	cdbm := &ArcCdbm{}
	cdbm.BasicCdbm = base.NewBasicCdbm()
	cdbm.t1 = NewCdbList()
	cdbm.b1 = NewCdbList()
	cdbm.t2 = NewCdbList()
	cdbm.b2 = NewCdbList()
	return cdbm
}

func (cdbm *ArcCdbm) Find(key string) (base.CacheDirectoryBlock, error) {
	cdb, err := cdbm.BasicCdbm.Find(key)
	if err == nil {
		acdb := cdb.(*ArcCdb)
		where := acdb.where
		if where == in_t1 {
			cdbm.t1.Remove(acdb.element)
			acdb.element = cdbm.t2.PushBack(acdb)
			acdb.where = in_t2
		} else if where == in_t2 {
			cdbm.t2.MoveToBack(acdb.element)
		} else {
			err = CacheMiss
		}
	}
	return cdb, err
}

func (cdbm *ArcCdbm) Remove(key string, f CacheCleanFunc) int {
	cdb, ok := cdbm.Hash[key]
	if ok {
		object := cdb.GetObject()
		oSize := object.Size()
		acdb := cdb.(*ArcCdb)
		where := acdb.where
		cdbm.Size -= oSize
		if f != nil {
			f(object)
		}
		if where == in_t1 {
			cdbm.t1.Remove(acdb.element)
		} else if  where == in_b1 {
			cdbm.target_t1 = min(cdbm.target_t1 + max(cdbm.b2.size/cdbm.b1.size, 1), cdbm.Size)
			cdbm.b1.Remove(acdb.element)
		} else  if where == in_t2 {
			cdbm.t2.Remove(acdb.element)
		} else if where == in_b2 {
			cdbm.target_t1 = max(cdbm.target_t1 - max(cdbm.b1.size/cdbm.b2.size, 1), 0)
			cdbm.b2.Remove(acdb.element)
		}
		delete(cdbm.Hash, key)
		return where
	}
	return not_in
}

func (cdbm *ArcCdbm) evict(f CacheCleanFunc) {
	var cdb base.CacheDirectoryBlock
	if cdbm.t1.size >= max(1, cdbm.target_t1) {
		cdb = cdbm.t1.RemoveLRU()
		cdb.(*ArcCdb).where = in_b1
		cdb.(*ArcCdb).element = cdbm.b1.PushBack(cdb)
	} else {
		cdb = cdbm.t2.RemoveLRU()
		cdb.(*ArcCdb).where = in_b2
		cdb.(*ArcCdb).element = cdbm.b1.PushBack(cdb)
	}
	object := cdb.GetObject()
	cdbm.Size -= object.Size()
	if f != nil {
		f(object)
	}
	return
}

func (cdbm *ArcCdbm) MakeSpace(objectSize, sizeLimit int, f CacheCleanFunc) (base.CacheDirectoryBlock, error) {
	if sizeLimit < objectSize {
		return nil, ObjectTooBig
	}
	
	//there is nothing 
	if len(cdbm.Hash) == 0 {
		return nil, nil
	}

	var cdb base.CacheDirectoryBlock
	for avail := sizeLimit - cdbm.Size; objectSize > avail; avail = sizeLimit - cdbm.Size {
		if cdbm.t1.size + cdbm.b1.size + objectSize <= cdbm.Size {
			if cdbm.b1.size == 0 {
				cdb = cdbm.t1.RemoveLRU()
			} else {
				cdb = cdbm.b1.RemoveLRU()
				cdbm.evict(f)
			}
		} else if cdbm.t1.size + cdbm.t2.size + cdbm.b1.size + cdbm.b2.size + objectSize >= cdbm.Size {
			if cdbm.t1.size + cdbm.t2.size + cdbm.b1.size + cdbm.b2.size + objectSize >= cdbm.Size * 2 {
				cdb = cdbm.b2.RemoveLRU()
			} else {
				cdb = newCacheDirectorBlock()
			}
			cdbm.evict(f)
		} else {
			cdb = newCacheDirectorBlock()
		}
	}
	cdb.(*ArcCdb).size = objectSize
	return cdb, nil
}

func (cdbm *ArcCdbm) Replace(key string, object CacheObject, sizeLimit int, f CacheCleanFunc) error {
	where := cdbm.Remove(key, f)
	oSize := object.Size()
	cdb, err := cdbm.MakeSpace(oSize, sizeLimit, f)
	if err != nil {
		return err
	}
	if cdb == nil {
		cdb = newCacheDirectorBlock()
	}
	acdb := cdb.(*ArcCdb)
	if where == not_in {
		acdb.where = in_t1
		acdb.element = cdbm.t1.PushBack(cdb)
	} else {
		acdb.where = in_t2
		acdb.element = cdbm.t2.PushBack(cdb)
	}
	cdb.SetKey(key)
	cdb.SetObject(object)
	cdbm.Size += oSize
	cdbm.Hash[key] = cdb
	return nil
}

func (cdbm *ArcCdbm) Check() {
	for key, cdb := range(cdbm.Hash) {
		if cdb.GetKey() != key {
			panic("keys don't match " + cdb.GetKey() + "!=" + key)
		}
		object := cdb.GetObject()
		if (cdb.(*ArcCdb).where == in_b1 || cdb.(*ArcCdb).where == in_b2) && object != nil {
			panic("cdb pointer should be nil")
		}
		if (cdb.(*ArcCdb).where == in_t1 || cdb.(*ArcCdb).where == in_t2) && object == nil {
			panic("cdb pointer should not be nil")
		}
	}
}