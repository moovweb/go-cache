package test

import "testing"
import "github.com/moovweb/go-cache"
import "github.com/moovweb/go-cache/arc"
import "go-cache/lru"
import "github.com/moovweb/go-cache/base"
import "strings"
import "io/ioutil"
import "sync"

type StringObject struct {
	s string
}

func (o *StringObject) Size() int {
	return len(o.s)
}

const cacheSize = 20 * 40
const concurrency = 20

func TestARC(t *testing.T) {
	countCleaned := 0
	countAdded := 0
	countMiss := 0

	c := arc.NewSafeARCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	c.SetCleanFunc(func(obj cache.CacheObject) error {
		countCleaned += obj.Size()
		return nil
	})
	countAccess := len(lines)
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i++ {
				val, err := c.Get(lines[i])
				if err == cache.CacheMiss {
					countAdded += len(lines[i])
					c.Set(lines[i], &StringObject{s: lines[i]})
					countMiss += 1
				} else if val.(*StringObject).s != lines[i] {
					t.Errorf("key does not match the value")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.Check()

	if countCleaned+c.GetUsage() != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, c.GetUsage())
	}

	for key, obj := range c.Collect() {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	println("cache hit rate/usage/# in cache:", c.GetHitRate(), c.GetUsage(), len(c.Collect()))
}

func TestLRU(t *testing.T) {
	countCleaned := 0
	countAdded := 0
	countMiss := 0

	c := lru.NewSafeLRUCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	c.SetCleanFunc(func(obj cache.CacheObject) error {
		countCleaned += obj.Size()
		return nil
	})
	countAccess := len(lines)
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i++ {
				val, err := c.Get(lines[i])
				if err == cache.CacheMiss {
					countAdded += len(lines[i])
					c.Set(lines[i], &StringObject{s: lines[i]})
					countMiss += 1
				} else if val.(*StringObject).s != lines[i] {
					t.Errorf("key does not match the value")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.Check()

	if countCleaned+c.GetUsage() != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, c.GetUsage())
	}

	for key, obj := range c.Collect() {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	println("cache hit rate/usage/# in cache:", c.GetHitRate(), c.GetUsage(), len(c.Collect()))
}

func TestRandom(t *testing.T) {
	countCleaned := 0
	countAdded := 0
	countMiss := 0

	c := base.NewSafeRRCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	c.SetCleanFunc(func(obj cache.CacheObject) error {
		countCleaned += obj.Size()
		return nil
	})
	countAccess := len(lines)
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i++ {
				val, err := c.Get(lines[i])
				if err == cache.CacheMiss {
					countAdded += len(lines[i])
					c.Set(lines[i], &StringObject{s: lines[i]})
					countMiss += 1
				} else if val.(*StringObject).s != lines[i] {
					t.Errorf("key does not match the value")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.Check()

	if countCleaned+c.GetUsage() != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, c.GetUsage())
	}

	for key, obj := range c.Collect() {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	println("cache hit rate/usage/# in cache:", c.GetHitRate(), c.GetUsage(), len(c.Collect()))
}
