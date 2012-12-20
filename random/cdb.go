package rrc

import (
	"math/rand"
	"go-cache/base"
)

type CDBList struct {
	cdbs []base.CacheDirectoryBlock
}

func newCacheDirectorBlock() *base.BaseCdb {
	cdb := &base.BaseCdb{}
	return cdb
}

func newCdbList(size int) *CDBList {
	cdbl := &CDBList{}
	cdbl.cdbs = make([]base.CacheDirectoryBlock, 0, size)
	return cdbl
}

func (cdbl *CDBList) Add(cdb base.CacheDirectoryBlock) {
	cdbl.cdbs = append(cdbl.cdbs, cdb)
}

func (cdbl *CDBList) Len() int {
	return len(cdbl.cdbs)
}

func (cdbl *CDBList) Select() base.CacheDirectoryBlock {
	num := cdbl.Len()
	if num > 0 {
		return cdbl.cdbs[rand.Intn(num)]
	}
	return nil
}
