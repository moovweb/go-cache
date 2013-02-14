package base

import . "go-cache"
import "math/rand"

type CdbManager interface {
	Find(string) (CacheDirectoryBlock, error)
	Replace(string, CacheObject, int, CacheCleanFunc) error
	Remove(key string, f CacheCleanFunc) error
	Collect() map[string]CacheObject
	Check()
	GetUsage() int
	Reset(CacheCleanFunc)
}

type BasicCdbm struct {
	//hash table for searching cache entries
	Hash map[string]CacheDirectoryBlock
	Size int
}

func NewBasicCdbm() *BasicCdbm {
	cdbm := &BasicCdbm{}
	cdbm.Hash = make(map[string]CacheDirectoryBlock)
	return cdbm
}

func (cdbm *BasicCdbm) Find(key string) (CacheDirectoryBlock, error) {
	cdb, ok := cdbm.Hash[key]
	if !ok {
		return nil, CacheMiss
	}
	return cdb, nil
}

func (cdbm *BasicCdbm) MakeSpace(objectSize, sizeLimit int, f CacheCleanFunc) (CacheDirectoryBlock, error) {
	if sizeLimit < objectSize {
		return nil, ObjectTooBig
	}

	//there is nothing 
	if len(cdbm.Hash) == 0 {
		return nil, nil
	}

	var repl CacheDirectoryBlock

	for avail := sizeLimit - cdbm.Size; objectSize > avail; avail = sizeLimit - cdbm.Size {
		cdbs := make([]CacheDirectoryBlock, 0, len(cdbm.Hash))
		for _, val := range cdbm.Hash {
			cdbs = append(cdbs, val)
		}
		num := len(cdbs)
		if num <= 0 {
			panic("cdbs is empty with nonzero size")
		}
		repl = cdbs[rand.Intn(num)]
		evObject := repl.GetObject()
		cdbm.Size -= evObject.Size()
		if f != nil && evObject != nil {
			f(evObject)
		}
		delete(cdbm.Hash, repl.GetKey())
	}
	return repl, nil
}

func (cdbm *BasicCdbm) Remove(key string, f CacheCleanFunc) error {
	cdb, ok := cdbm.Hash[key]
	if ok {
		object := cdb.GetObject()
		cdbm.Size -= object.Size()
		if f != nil {
			f(object)
		}
		delete(cdbm.Hash, key)
	}
	return nil
}

func (cdbm *BasicCdbm) Replace(key string, object CacheObject, sizeLimit int, f CacheCleanFunc) error {
	cdbm.Remove(key, f)
	oSize := object.Size()
	cdb, err := cdbm.MakeSpace(oSize, sizeLimit, f)
	if err != nil {
		return err
	}
	if cdb == nil {
		cdb = NewBasicCdb()
	}
	cdb.SetKey(key)
	cdb.SetObject(object)
	cdbm.Size += oSize
	cdbm.Hash[key] = cdb
	return nil
}

func (cdbm *BasicCdbm) Collect() map[string]CacheObject {
	m := make(map[string]CacheObject)
	for key, cdb := range cdbm.Hash {
		m[key] = cdb.GetObject()
	}
	return m
}

func (cdbm *BasicCdbm) Check() {
	for key, cdb := range cdbm.Hash {
		if cdb.GetKey() != key {
			panic("keys don't match")
		}
	}
	if cdbm.GetUsage() > cdbm.Size {
		panic("Cache usage exceeds the limit")
	}
}

func (cdbm *BasicCdbm) GetUsage() int {
	return cdbm.Size
}

func (cdbm *BasicCdbm) Reset(f CacheCleanFunc) {
	for key, _ := range cdbm.Hash {
		cdbm.Remove(key, f)
	}
}
