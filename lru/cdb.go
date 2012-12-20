package lru

import (
	"container/list"
	"go-cache/base"
)

type CDBList struct {
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

func newCdbList() *CDBList {
	cdbl := &CDBList{}
	cdbl.cdbs = list.New()
	return cdbl
}

func (cdbl *CDBList) RemoveLRU() base.CacheDirectoryBlock {
	lru := cdbl.cdbs.Front()
	if lru == nil {
		return nil
	}
	cdb := cdbl.cdbs.Remove(lru).(base.CacheDirectoryBlock)
	return cdb
}

func (cdbl *CDBList) InsertMRU(cdb base.CacheDirectoryBlock) {
	cdb.(*ElementCdb).element = cdbl.cdbs.PushBack(cdb)
}

func (cdbl *CDBList) SetMRU(cdb base.CacheDirectoryBlock) {
	cdbl.cdbs.MoveToBack(cdb.(*ElementCdb).element)
}

func (cdbl *CDBList) RemoveIt(cdb base.CacheDirectoryBlock) {
	cdbl.cdbs.Remove(cdb.(*ElementCdb).element)
}

func (cdbl *CDBList) Len() int {
	return cdbl.cdbs.Len()
}
