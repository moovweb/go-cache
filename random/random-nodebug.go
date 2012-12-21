// +build !cache_debug

package rrc

import . "go-cache"

func (c *RRCache) Get(key string) (object CacheObject, err error) {
	if len(key) == 0 {
		return nil, EmptyKey
	}
	cdb, err := c.get(key)
	return c.GetOrFetch(key, cdb, err)
}