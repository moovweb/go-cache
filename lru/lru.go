package lru

import (
	"container/list"
)

import . "go-cache"

type LRUCache struct {
	maxItems int
	lru      *list.List
	index    map[string]*list.Element
	
	//the hook which is called in case of cache miss
	fetchFunc CacheFetchFunc

	//clean func
	cleanFunc CacheCleanFunc
}

type keyValue struct {
	key   string
	value interface{}
}

func NewLRUCache(maxItems int) *LRUCache {
	cache := &LRUCache{
		maxItems: maxItems,
		lru:      list.New(),
		index:    make(map[string]*list.Element, maxItems),
	}
	return cache
}

func (c *LRUCache) SetFetchFunc(f CacheFetchFunc) {
	c.fetchFunc = f
}

func (c *LRUCache) SetCleanFunc(f CacheCleanFunc) {
	c.cleanFunc = f
}

func (c *LRUCache) Set(key string, object CacheObject) {
	el, ok := c.index[key]
	if ok {
		// Element exists so just move it to the back and update value
		kv := el.Value.(*keyValue)
		kv.value = object
		c.lru.MoveToBack(el)
	} else {
		if len(c.index) >= c.maxItems {
			// Cache is full so remove an existing key/value
			el := c.lru.Front()
			kv := el.Value.(*keyValue)
			c.delete(kv)
			// Reuse list element
			kv.key = key
			kv.value = object
			c.lru.MoveToBack(el)
			c.index[key] = el
		} else {
			// Cache is not full and this is a new key
			kv := &keyValue{key, object}
			el := c.lru.PushBack(kv)
			c.index[key] = el
		}
	}
}

func (c *LRUCache) Get(key string) (CacheObject, error) {
	el, ok := c.index[key]
	if !ok {
		if c.fetchFunc == nil {
			return nil, CacheMiss
		}
		object, err := c.fetchFunc(key)
		c.Set(key, object)
		if err != nil {
			return object, err
		}
		return object, CacheMiss
	}
	c.lru.MoveToBack(el)
	kv := el.Value.(*keyValue)
	return kv.value, nil
}

func (c *LRUCache) delete(kv *keyValue) {
	if c.cleanFunc != nil {
		c.cleanFunc(kv.value)
	}
	delete(c.index, kv.key)
}

func (c *LRUCache) GetAllObjects() map[string]CacheObject {
	all := make(map[string]CacheObject)
	for key, value := range c.index {
		all[key] = value.Value.(*keyValue).value
	}
	return all
}


func (c *LRUCache) CheckCache() {
	for key, value := range(c.index) {
		if value.Value.(*keyValue).key != key {
			panic("keys don't match")
		}
	}
}
