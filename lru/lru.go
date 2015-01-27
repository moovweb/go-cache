package lru

import "github.com/moovweb/go-cache/base"

func NewLRUCache(size int) *base.BaseCache {
	lru := newCdbm()
	c := base.NewBaseCache(size, lru)
	return c
}

func NewSafeLRUCache(size int) *base.BaseCache {
	lru := newCdbm()
	c := base.NewSafeBaseCache(size, lru)
	return c
}
