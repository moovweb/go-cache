package rrc

import (
	"sync"
	"math/rand"
)

type CDBList struct {
	cdbs []*cacheEntry
	mutex sync.Mutex
}

type cacheDirectoryBlock struct {
	index int
	pointer *cacheEntry
	key string
}

func newCacheDirectorBlock() *cacheDirectoryBlock {
	cdb := &cacheDirectoryBlock{}
	return cdb
}

func newCdbList(size int) *CDBList {
	cdbl := &CDBList{}
	cdbl.cdbs = make([]*cacheEntry, 0, size)
	return cdbl
}

func (cdbl *CDBList) Add(cdb *cacheDirectoryBlock) {
	index := len(cdbl.cdbs)
	cdbl.cdbs = append(cdbl.cdbs, cdb.pointer)
	cdb.index = index
}

func (cdbl *CDBList) Len() int {
	return len(cdbl.cdbs)
}

func (cdbl *CDBList) RandomSelection() *cacheDirectoryBlock {
	num := cdbl.Len()
	if num > 0 {
		return cdbl.cdbs[rand.Intn(num)].cdb
	}
	return nil
}
