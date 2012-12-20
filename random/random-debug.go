// +build cache_debug

package rrc

import . "go-cache"
import "time"

func (c *RRCache) Get(key string) (object CacheObject, err error) {
	start := time.Now()
	tmp, err := c.get(key)
	entry := tmp.GetEntry()
	t := time.Since(start)
	c.AddAccessTime(t.Nanoseconds())
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