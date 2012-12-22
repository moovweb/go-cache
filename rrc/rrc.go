package rrc

import . "go-cache"
import "go-cache/base"

func NewRRCache(size int) Cache {
	cdbm := newCdbm()
	c := base.NewBaseCache(size, cdbm)
	return c
}

func NewSafeRRCache(size int) Cache {
	cdbm := newCdbm()
	c := base.NewSafeBaseCache(size, cdbm)
	return c
}