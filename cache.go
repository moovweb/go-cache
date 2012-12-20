package cache

import "errors"

type CacheObject interface {
	Size() int
}

type CacheFetchFunc func(string) (CacheObject, error)
type CacheCleanFunc func(CacheObject) error

type Cache interface {
	Set(string, CacheObject)
	Get(string) (CacheObject, error)
	GetAllObjects() map[string]CacheObject
	SetFetchFunc(CacheFetchFunc)
	SetCleanFunc(CacheCleanFunc)
	GetHitRate() int
	Resize(size int)
}

var CacheMiss = errors.New("miss")