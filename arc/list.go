package arc

import "container/list"
//import . "go-cache"
import "go-cache/base"

type CdbList struct {
	l *list.List
	size int
}


func (cl *CdbList) PushBack(cdb base.CacheDirectoryBlock) *list.Element {
	e := cl.l.PushBack(cdb)
	cl.size += cdb.(*ArcCdb).size
	return e
}


func (cl *CdbList) Remove(e *list.Element) base.CacheDirectoryBlock {
	cdb := cl.l.Remove(e).(base.CacheDirectoryBlock)
	cl.size -= cdb.(*ArcCdb).size
	return cdb
}

func (cl *CdbList) MoveToBack(e *list.Element) {
	cl.l.MoveToBack(e)
}

func (cl *CdbList) RemoveLRU() base.CacheDirectoryBlock {
	lru := cl.l.Front()
	cdb := cl.l.Remove(lru).(base.CacheDirectoryBlock)
	cl.size -= cdb.(*ArcCdb).size
	return cdb
}

func NewCdbList() *CdbList {
	cl := &CdbList{}
	cl.l = list.New()
	return cl
}

func (cl *CdbList) Len() int {
	return cl.l.Len()
}