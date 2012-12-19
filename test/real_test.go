package test

import "testing"
import "go-cache"
import "go-cache/arc"
import "go-cache/lru"
import "go-cache/random"
import "strings"
import "io/ioutil"

func TestARC(t *testing.T) {
	cacheSize := 25
	countCleaned := 0
	countAdded := 0

	c := arc.NewArcCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	c.SetCleanFunc(func (obj cache.CacheObject) error {
		countCleaned += 1
		return nil
	})
	c.SetFetchFunc(func (key string) (cache.CacheObject, error) {
		countAdded += 1
		return key, nil
	})

	countMiss := 0
	countAccess := len(lines)
	for i := 0; i < countAccess; i ++ {
		_, err := c.Get(lines[i])
		if err == cache.CacheMiss {
			countMiss += 1
		}
	}
	println("cache hit rate:", (100*(countAccess-countMiss))/countAccess)

	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(string) {
			t.Errorf("key does not match the cached value")
		}
	}
}

func TestLRU(t *testing.T) {
	cacheSize := 25
	countCleaned := 0
	countAdded := 0

	c := lru.NewLRUCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	c.SetCleanFunc(func (obj cache.CacheObject) error {
		countCleaned += 1
		return nil
	})
	c.SetFetchFunc(func (key string) (cache.CacheObject, error) {
		countAdded += 1
		return key, nil
	})

	countMiss := 0
	countAccess := len(lines)
	for i := 0; i < countAccess; i ++ {
		_, err := c.Get(lines[i])
		if err == cache.CacheMiss {
			countMiss += 1
		}
	}
	println("cache hit rate:", (100*(countAccess-countMiss))/countAccess)

	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(string) {
			t.Errorf("key does not match the cached value")
		}
	}
}

func TestRRC(t *testing.T) {
	cacheSize := 25
	countCleaned := 0
	countAdded := 0

	c := rrc.NewRRCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	c.SetCleanFunc(func (obj cache.CacheObject) error {
		countCleaned += 1
		return nil
	})
	c.SetFetchFunc(func (key string) (cache.CacheObject, error) {
		countAdded += 1
		return key, nil
	})

	countMiss := 0
	countAccess := len(lines)
	for i := 0; i < countAccess; i ++ {
		_, err := c.Get(lines[i])
		if err == cache.CacheMiss {
			countMiss += 1
		}
	}
	println("cache hit rate:", (100*(countAccess-countMiss))/countAccess)

	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(string) {
			t.Errorf("key does not match the cached value")
		}
	}
}