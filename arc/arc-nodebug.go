// +build !cache_debug

package arc

import . "go-cache"

func (c *ARCache) Get(key string) (object CacheObject, err error) {
	if len(key) == 0 {
		return nil, EmptyKey
	}
	cdb, err := c.get(key)
	return c.GetOrFetch(key, cdb, err)
}
