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
	key string
}

func newCacheDirectorBlock() *cacheDirectoryBlock {
	cdb := &cacheDirectoryBlock{}
	return cdb
}

func newCdbList() *CDBList {
	cdbl := &CDBList{}
	cdbl.cdbs = list.New()
	return cdbl
}

func (cdbl *CDBList) RemoveLRU() *cacheDirectoryBlock {
	lru := cdbl.cdbs.Front()
	if lru == nil {
		return nil
	}
	cdb := cdbl.cdbs.Remove(lru).(*cacheDirectoryBlock)
	return cdb
}

func (cdbl *CDBList) InsertMRU(cdb *cacheDirectoryBlock) {
	cdb.element = cdbl.cdbs.PushBack(cdb)
}

func (cdbl *CDBList) SetMRU(cdb *cacheDirectoryBlock) {
	cdbl.cdbs.MoveToBack(cdb.element)
}

func (cdbl *CDBList) RemoveIt(cdb *cacheDirectoryBlock) {
	cdb = cdbl.cdbs.Remove(cdb.element).(*cacheDirectoryBlock)
}

func (cdbl *CDBList) Len() int {
	return cdbl.cdbs.Len()
}
