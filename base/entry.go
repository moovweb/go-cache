package base

import . "go-cache"
import "sync"

type CacheEntry struct {
	object CacheObject
	mutex sync.RWMutex
}

func NewCacheEntry() *CacheEntry {
	return &CacheEntry{}
}

func (entry *CacheEntry) GetObject() CacheObject {
	entry.mutex.RLock()
	defer entry.mutex.RUnlock()
	return entry.object
}

func (entry *CacheEntry) SetObject(obj CacheObject, f CacheCleanFunc) {
	entry.mutex.Lock()
	defer entry.mutex.Unlock()
	if entry.object != nil && f != nil {
		f(entry.object)
	}
	entry.object = obj
}