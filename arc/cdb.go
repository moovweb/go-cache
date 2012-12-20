package arc

import (
	"container/list"
	"go-cache/base"
)

const (
	in_t1 = iota
	in_b1
	in_t2
	in_b2
)

type CDBList struct {
	cdbs *list.List
}

type ArcCdb struct {
	element *list.Element
	where int
	base.BaseCdb
}

func newCacheDirectorBlock() *ArcCdb {
	cdb := &ArcCdb{}
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
	cdb.(*ArcCdb).element = cdbl.cdbs.PushBack(cdb)
}

func (cdbl *CDBList) SetMRU(cdb base.CacheDirectoryBlock) {
	cdbl.cdbs.MoveToBack(cdb.(*ArcCdb).element)
}

func (cdbl *CDBList) RemoveIt(cdb base.CacheDirectoryBlock) {
	cdbl.cdbs.Remove(cdb.(*ArcCdb).element)
}

func (cdbl *CDBList) Len() int {
	return cdbl.cdbs.Len()
}
