package base

import . "go-cache"
import "sync"

type CacheEntry struct {
	object CacheObject
	dirty bool
	mutex sync.Mutex
	cond *sync.Cond
	isGoroutineSafe bool
	waits int
}

func NewCacheEntry() *CacheEntry {
	entry := &CacheEntry{isGoroutineSafe: false}
	if entry.isGoroutineSafe {
		entry.cond = sync.NewCond(&entry.mutex)
	}
	return entry
}

func NewSafeCacheEntry() *CacheEntry {
	entry := &CacheEntry{isGoroutineSafe: true}
	if entry.isGoroutineSafe {
		entry.cond = sync.NewCond(&entry.mutex)
	}
	return entry
}

func (entry *CacheEntry) GetObject() CacheObject {
	if entry.isGoroutineSafe {
		entry.mutex.Lock()
		defer entry.mutex.Unlock()
		for(entry.dirty) {
			entry.waits += 1
			entry.cond.Wait()
			entry.waits -= 1
		}
	}
	return entry.object
}

func (entry *CacheEntry) SetObject(obj CacheObject, f CacheCleanFunc) {
	if entry.isGoroutineSafe {
		entry.mutex.Lock()
		defer entry.mutex.Unlock()
	}
	if entry.object != nil && f != nil {
		f(entry.object)
	}
	entry.object = obj
	entry.dirty = false
	if entry.isGoroutineSafe {
		for i := entry.waits; i > 0; i-- {
			entry.cond.Signal()
		}
	}
}

func (entry *CacheEntry) SetDirty() {
	if entry.isGoroutineSafe {
		entry.mutex.Lock()
		defer entry.mutex.Unlock()
	}
	entry.dirty = true
}