package rrc

import (
	"math/rand"
	. "go-cache"
	"go-cache/base"
)

type BasicCdbm struct {
	//hash table for searching cache entries
	hash map[string]base.CacheDirectoryBlock
	size int
}

func newCacheDirectorBlock() *base.BaseCdb {
	cdb := &base.BaseCdb{}
	return cdb
}

func newCdbm() *BasicCdbm {
	cdbm := &BasicCdbm{}
	cdbm.hash = make(map[string]base.CacheDirectoryBlock)
	return cdbm
}

func (cdbm *BasicCdbm) Find(key string) (base.CacheDirectoryBlock, error) {
	cdb, ok := cdbm.hash[key]
	if !ok {
		return nil, CacheMiss
	}
	return cdb, nil
}

func (cdbm *BasicCdbm) MakeSpace(objectSize, sizeLimit int, f CacheCleanFunc) (base.CacheDirectoryBlock, error) {
	if sizeLimit < objectSize {
		return nil, ObjectTooBig
	}
	
	//there is nothing 
	if len(cdbm.hash) == 0 {
		return nil, nil
	}

	var repl base.CacheDirectoryBlock
	
	for avail := sizeLimit - cdbm.size; objectSize > avail; avail = sizeLimit - cdbm.size {
		cdbs := make([]base.CacheDirectoryBlock, 0, len(cdbm.hash))
		for _, val := range(cdbm.hash) {
			cdbs = append(cdbs, val)
		}
		num := len(cdbs)
		if num <= 0 {
			panic("cdbs is empty with nonzero size")
		}
		repl = cdbs[rand.Intn(num)]
		evObject := repl.GetObject()
		cdbm.size -= evObject.Size()
		f(evObject)
		println("evict")
		delete(cdbm.hash, repl.GetKey())
	}
	return repl, nil
}

func (cdbm *BasicCdbm) Replace(key string, object CacheObject, sizeLimit int, f CacheCleanFunc) error {
	cdb, ok := cdbm.hash[key]
	if !ok {
		oSize := object.Size()
		println("oSize:", oSize, "size:", )
		repl, err := cdbm.MakeSpace(oSize, sizeLimit, f)
		if err != nil {
			return err
		}
		
		if repl == nil {
			cdb = newCacheDirectorBlock()
			println("create new")
		} else {
			cdb = repl
		}
	}
	cdb.SetKey(key)
	cdb.SetObject(object)
	cdbm.size += object.Size()
	cdbm.hash[key] = cdb
	return nil
}

func (cdbm *BasicCdbm) Collect() map[string]CacheObject {
	m := make(map[string]CacheObject)
	for key, cdb := range(cdbm.hash) {
		m[key] = cdb.GetObject()
	}
	return m
}

func (cdbm *BasicCdbm) Check() {
	for key, cdb := range(cdbm.hash) {
		if cdb.GetKey() != key {
			panic("keys don't match")
		}
	}
}

