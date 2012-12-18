package cache

type CacheObject interface {}

type CacheFetchFunc func(string) (CacheObject, error)
type CacheCleanFunc func(CacheObject) error

type Cache interface {
	Set(string, CacheObject)
	Get(string) (CacheObject, error)
	SetFetchFunc(CacheFetchFunc)
	SetCleanFunc(CacheCleanFunc)
}

