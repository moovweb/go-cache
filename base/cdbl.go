package base

import . "go-cache"

type CdbManager interface {
	Find(string) (CacheDirectoryBlock, error)
	Replace(string, CacheObject, int, CacheCleanFunc) error
	Collect() map[string]CacheObject
	Check()
	GetUsage() int
}