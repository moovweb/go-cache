// +build !cache_debug

package rrc

import . "go-cache"

func (c *RRCache) Get(key string) (object CacheObject, err error) {
	tmp, err := c.get(key)
	entry := tmp.GetEntry()
	if err == CacheMiss {
		var err1 error
		object, err1 = c.FetchFunc(key)
		entry.SetObject(object, c.CleanFunc)
		if err1 != nil {
			err = err1
		}
	} else {
		object = entry.GetObject()
	}
	return
}