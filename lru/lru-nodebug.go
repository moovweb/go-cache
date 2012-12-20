// +build !cache_debug

package lru

import . "go-cache"

func (c *LRUCache) Get(key string) (object CacheObject, err error) {
	tmp, err := c.get(key)
	entry := tmp.GetEntry()
	c.Accesses += 1
	if err == CacheMiss {
		var err1 error
		object, err1 = c.FetchFunc(key)
		entry.SetObject(object, c.CleanFunc)
		if err1 != nil {
			err = err1
		}
	} else {
		c.Hits += 1
		object = entry.GetObject()
	}
	return
}