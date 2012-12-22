package rrc

import "go-cache/base"

type RRCache struct {
	*base.BaseCache
}

func NewRRCache(size int) *RRCache {
	c:= &RRCache{}
	cdbm := newCdbm()
	c.BaseCache = base.NewBaseCache(size, cdbm)
	return c
}

func NewSafeRRCache(size int) *RRCache {
	c:= &RRCache{}
	cdbm := newCdbm()
	c.BaseCache = base.NewSafeBaseCache(size, cdbm)
	return c
}