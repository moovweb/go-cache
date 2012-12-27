package arc

import "go-cache/base"

func NewARCache(size int) *base.BaseCache {
	arc := newCdbm()
	c := base.NewBaseCache(size, arc)
	return c
}

func NewSafeARCache(size int) *base.BaseCache {
	arc := newCdbm()
	c := base.NewSafeBaseCache(size, arc)
	return c
}