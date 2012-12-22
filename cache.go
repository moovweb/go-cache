package cache

import "errors"

type CacheObject interface {
	Size() int
}

type CacheCleanFunc func(CacheObject) error

type Cache interface {
	Set(string, CacheObject)
	Get(string) (CacheObject, error)
	Collect() map[string]CacheObject
	SetCleanFunc(CacheCleanFunc)
}

var CacheMiss = errors.New("miss")
var EmptyKey = errors.New("empty key")
var ObjectTooBig = errors.New("object too big to fit in the cache")