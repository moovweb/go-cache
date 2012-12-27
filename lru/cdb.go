package lru

import (
	"container/list"
	. "go-cache"
	"go-cache/base"
)

type LruCdbm struct {
	*base.BasicCdbm
	cdbs *list.List
}

type ElementCdb struct {
	element *list.Element
	base.BaseCdb
}

func newCacheDirectorBlock() *ElementCdb {
	cdb := &ElementCdb{}
	return cdb
}

func newCdbm() *LruCdbm {
	cdbm := &LruCdbm{}
	cdbm.BasicCdbm = base.NewBasicCdbm()
	cdbm.cdbs = list.New()
	return cdbm
}

func (cdbm *LruCdbm) Find(key string) (base.CacheDirectoryBlock, error) {
	cdb, err := cdbm.BasicCdbm.Find(key)
	if err == nil {
		cdbm.cdbs.MoveToBack(cdb.(*ElementCdb).element)
	}
	return cdb, err
}

func (cdbm *LruCdbm) Remove(key string, f CacheCleanFunc) error {
	cdb, ok := cdbm.Hash[key]
	if ok {
		object := cdb.GetObject()
		cdbm.Size -= object.Size()
		if f != nil {
			f(object)
		}
		cdbm.cdbs.Remove(cdb.(*ElementCdb).element)
		delete(cdbm.Hash, key)
	}
	return nil
}

func (cdbm *LruCdbm) MakeSpace(objectSize, sizeLimit int, f CacheCleanFunc) (base.CacheDirectoryBlock, error) {
	if sizeLimit < objectSize {
		return nil, ObjectTooBig
	}
	
	//there is nothing 
	if len(cdbm.Hash) == 0 {
		return nil, nil
	}

	var repl base.CacheDirectoryBlock
	
	for avail := sizeLimit - cdbm.Size; objectSize > avail; avail = sizeLimit - cdbm.Size {
		lru := cdbm.cdbs.Front()
		if lru == nil {
			panic("cdbs is empty with nonzero size")
		}
		repl = cdbm.cdbs.Remove(lru).(base.CacheDirectoryBlock)
		evObject := repl.GetObject()
		cdbm.Size -= evObject.Size()
		if f != nil && evObject != nil {
			f(evObject)
		}
		delete(cdbm.Hash, repl.GetKey())
	}
	return repl, nil
}

func (cdbm *LruCdbm) Replace(key string, object CacheObject, sizeLimit int, f CacheCleanFunc) error {
	cdbm.Remove(key, f)
	oSize := object.Size()
	cdb, err := cdbm.MakeSpace(oSize, sizeLimit, f)
	if err != nil {
		return err
	}
	
	if cdb == nil {
		cdb = newCacheDirectorBlock()
	}
	cdb.SetKey(key)
	cdb.SetObject(object)
	cdbm.Size += object.Size()
	cdbm.Hash[key] = cdb
	cdb.(*ElementCdb).element = cdbm.cdbs.PushBack(cdb)
	return nil
}