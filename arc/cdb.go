package arc

import (
	"container/list"
	"sync"
)

const (
	in_t1 = iota
	in_b1
	in_t2
	in_b2
)

type CDBList struct {
	cdbs *list.List
	mutex sync.Mutex
}

type cacheDirectoryBlock struct {
	element *list.Element
	pointer *cacheEntry
	where int
}

func newCacheDirectorBlock() *cacheDirectoryBlock {
	cdb := &cacheDirectoryBlock{}
	return cdb
}

func (cdbl *CDBList) RemoveLRU() *cacheDirectoryBlock {
	lru := cdbl.cdbs.Front()
	if lru == nil {
		return nil
	}
	cdb := cdbl.cdbs.Remove(lru).(*cacheDirectoryBlock)
	if cdb != nil && cdb.pointer != nil && cdb.pointer.object != nil {
		cdb.pointer.object.Free()
		cdb.pointer.object = nil
	}
	return cdb
}

func (cdbl *CDBList) InsertMRU(cdb *cacheDirectoryBlock) {
	cdb.element = cdbl.cdbs.PushBack(cdb)
}

func (cdbl *CDBList) SetMRU(cdb *cacheDirectoryBlock) {
	cdbl.cdbs.MoveToBack(cdb.element)
}

func (cdbl *CDBList) RemoveIt(cdb *cacheDirectoryBlock) {
	cdbl.cdbs.Remove(cdb.element)
}

func (cdbl *CDBList) Len() int {
	return cdbl.cdbs.Len()
}