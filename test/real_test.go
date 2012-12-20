package test

import "testing"
import "go-cache"
import "go-cache/arc"
import "go-cache/lru"
import "go-cache/random"
import "strings"
import "io/ioutil"

type StringObject struct {
	s string
}

func (o *StringObject) Size() int {
	return len(o.s)
}

func TestARC(t *testing.T) {
	cacheSize := 30
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
		return &StringObject{s:key}, nil
	})

	countMiss := 0
	countAccess := len(lines)
	for i := 0; i < countAccess; i ++ {
		_, err := c.Get(lines[i])
		if err == cache.CacheMiss {
			countMiss += 1
		}
	}
	println("cache hit rate:", c.GetHitRate())
	//println("avg get time:", c.Total/c.Count)

	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
}

func TestLRU(t *testing.T) {
	cacheSize := 30
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
		return &StringObject{s:key}, nil
	})

	countAccess := len(lines)
	count := 0
	for i := 0; i < countAccess; i ++ {
		c.Get(lines[i])
		count += 1
	}
	println("cache hit rate:", c.GetHitRate())
	//println("avg get time:", c.Total/c.Count)

	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Fatalf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
}

func TestRRC(t *testing.T) {
	cacheSize := 30
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
		return &StringObject{s:key}, nil
	})

	countAccess := len(lines)

	for i := 0; i < countAccess; i ++ {
		c.Get(lines[i])
	}
	println("cache hit rate:", c.GetHitRate())
	//println("avg get time:", c.Total/c.Count)

	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
}