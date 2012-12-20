package rrc

import (
	"math/rand"
)

type CDBList struct {
	cdbs []*cacheDirectoryBlock
}

type cacheDirectoryBlock struct {
	pointer *cacheEntry
	key string
}

func newCacheDirectorBlock() *cacheDirectoryBlock {
	cdb := &cacheDirectoryBlock{}
	return cdb
}

func newCdbList(size int) *CDBList {
	cdbl := &CDBList{}
	cdbl.cdbs = make([]*cacheDirectoryBlock, 0, size)
	return cdbl
}

func (cdbl *CDBList) Add(cdb *cacheDirectoryBlock) {
	cdbl.cdbs = append(cdbl.cdbs, cdb)
}

func (cdbl *CDBList) Len() int {
	return len(cdbl.cdbs)
}

func (cdbl *CDBList) RandomSelection() *cacheDirectoryBlock {
	num := cdbl.Len()
	if num > 0 {
		return cdbl.cdbs[rand.Intn(num)]
	}
	return nil
}
