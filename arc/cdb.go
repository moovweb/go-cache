package arc

import (
	"container/list"
	"errors"

	. "go-cache"
	"go-cache/base"
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
	t1        *CdbList
	b1        *CdbList
	t2        *CdbList
	b2        *CdbList
	target_t1 int
}

type ArcCdb struct {
	element *list.Element
	where   int
	base.BaseCdb
	size int
	v    string
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

func (cdbm *ArcCdbm) remove(key string, f CacheCleanFunc) int {
	cdb, ok := cdbm.Hash[key]
	if ok {
		acdb := cdb.(*ArcCdb)
		where := acdb.where

		if where == in_t1 {
			object := cdb.GetObject()
			cdbm.t1.Remove(acdb.element)
			cdbm.Size -= object.Size()
			if f != nil {
				f(object)
			}
		} else if where == in_b1 {
			cdbm.target_t1 = min(cdbm.target_t1+max(cdbm.b2.size/cdbm.b1.size, 1), cdbm.Size)
			cdbm.b1.Remove(acdb.element)
		} else if where == in_t2 {
			object := cdb.GetObject()
			cdbm.t2.Remove(acdb.element)
			cdbm.Size -= object.Size()
			if f != nil {
				f(object)
			}
		} else if where == in_b2 {
			cdbm.target_t1 = max(cdbm.target_t1-max(cdbm.b1.size/cdbm.b2.size, 1), 0)
			cdbm.b2.Remove(acdb.element)
		}
		delete(cdbm.Hash, key)
		return where
	}
	return not_in
}

func (cbdm *ArcCdbm) Remove(key string, f CacheCleanFunc) error {
	where := cbdm.remove(key, f)
	if where == not_in {
		return errors.New("Tried to remove slug that doesn't exist: " + key)
	}
	return nil
}

func (cdbm *ArcCdbm) evict(f CacheCleanFunc) {
	var cdb base.CacheDirectoryBlock
	if cdbm.t1.size >= max(1, cdbm.target_t1) {
		cdb = cdbm.t1.RemoveLRU()
		acdb := cdb.(*ArcCdb)
		acdb.where = in_b1
		acdb.element = cdbm.b1.PushBack(acdb)
	} else {
		cdb = cdbm.t2.RemoveLRU()
		acdb := cdb.(*ArcCdb)
		acdb.where = in_b2
		acdb.element = cdbm.b2.PushBack(acdb)
	}
	object := cdb.GetObject()
	cdbm.Size -= object.Size()
	if f != nil {
		f(object)
	}
	cdb.SetObject(nil)
	return
}

func (cdbm *ArcCdbm) MakeSpace(objectSize, sizeLimit int, f CacheCleanFunc) (base.CacheDirectoryBlock, error) {
	if sizeLimit < objectSize {
		return nil, ObjectTooBig
	}

	var cdb base.CacheDirectoryBlock
	for avail := sizeLimit - cdbm.Size; objectSize > avail; avail = sizeLimit - cdbm.Size {
		if cdbm.t1.size+cdbm.b1.size+objectSize >= sizeLimit {
			if cdbm.b1.size > 0 {
				cdb = cdbm.b1.RemoveLRU()
				cdbm.evict(f)
			} else {
				cdb = cdbm.t1.RemoveLRU()
				object := cdb.GetObject()
				cdbm.Size -= object.Size()
				if f != nil {
					f(object)
				}
			}
			delete(cdbm.Hash, cdb.GetKey())
		} else if cdbm.t1.size+cdbm.t2.size+cdbm.b1.size+cdbm.b2.size+objectSize >= sizeLimit {
			if cdbm.t1.size+cdbm.t2.size+cdbm.b1.size+cdbm.b2.size+objectSize >= sizeLimit*2 {
				cdb = cdbm.b2.RemoveLRU()
				delete(cdbm.Hash, cdb.GetKey())
			} else {
				cdb = newCacheDirectorBlock()
			}
			cdbm.evict(f)
		} else {
			cdb = newCacheDirectorBlock()
		}
	}
	if cdb == nil {
		cdb = newCacheDirectorBlock()
	}
	return cdb, nil
}

func (cdbm *ArcCdbm) Replace(key string, object CacheObject, sizeLimit int, f CacheCleanFunc) error {
	where := cdbm.remove(key, f)
	oSize := object.Size()
	cdb, err := cdbm.MakeSpace(oSize, sizeLimit, f)
	if err != nil {
		return err
	}

	acdb := cdb.(*ArcCdb)
	acdb.size = oSize
	acdb.v = key
	if where == not_in {
		acdb.where = in_t1
		acdb.element = cdbm.t1.PushBack(acdb)
	} else {
		acdb.where = in_t2
		acdb.element = cdbm.t2.PushBack(acdb)
	}

	cdb.SetKey(key)
	cdb.SetObject(object)
	cdbm.Size += oSize
	cdbm.Hash[key] = cdb
	return nil
}

func (cdbm *ArcCdbm) Check() {
	for key, cdb := range cdbm.Hash {
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
	if cdbm.GetUsage() > cdbm.Size {
		panic("Cache usage exceeds the limit")
	}
}

func (cdbm *ArcCdbm) Collect() map[string]CacheObject {
	m := make(map[string]CacheObject)
	for key, cdb := range cdbm.Hash {
		if cdb.(*ArcCdb).where == in_t1 || cdb.(*ArcCdb).where == in_t2 {
			m[key] = cdb.GetObject()
		}
	}
	return m
}

func (cdbm *ArcCdbm) Reset(f CacheCleanFunc) {
	for key, _ := range cdbm.Hash {
		cdbm.Remove(key, f)
	}
}
