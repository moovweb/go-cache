package cache

import "errors"

type CacheObject interface {
	Size() int
}

type CacheCleanFunc func(CacheObject) error

type Cache interface {
	Set(string, CacheObject) error
	Get(string) (CacheObject, error)
	Check()
	Collect() map[string]CacheObject
	SetCleanFunc(CacheCleanFunc)
	GetHitRate() int
	GetUsage() int
	Reset()
}

var CacheMiss = errors.New("miss")
var EmptyKey = errors.New("empty key")
var ObjectTooBig = errors.New("object too big to fit in the cache")