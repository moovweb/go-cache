package cache

import "errors"

type CacheObject interface {}

type CacheFetchFunc func(string) (CacheObject, error)
type CacheCleanFunc func(CacheObject) error

type Cache interface {
	Set(string, CacheObject)
	Get(string) (CacheObject, error)
	GetAllObjects() map[string]CacheObject
	SetFetchFunc(CacheFetchFunc)
	SetCleanFunc(CacheCleanFunc)
}

var CacheMiss = errors.New("miss")