package arc

type CacheObject interface {
	Free()
}

type Cache interface {
	Set(string, CacheObject)
	Get(string) (CacheObject, error)
	Delete(string) error
}

type CacheFetchFunc func(string) (CacheObject, error)