package test

import "testing"
import "go-cache"
import "go-cache/arc"
import "go-cache/lru"
import "go-cache/random"
import "strings"
import "io/ioutil"
import "sync"

type StringObject struct {
	s string
}

func (o *StringObject) Size() int {
	return len(o.s)
}

const cacheSize = 20
const concurrency = 20

func TestARC(t *testing.T) {
	countCleaned := 0
	countAdded := 0

	c := arc.NewSafeArcCache(cacheSize)
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
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i ++ {
				c.Get(lines[i])
			}
			wg.Done()
		}()
	}
	wg.Wait()
	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	c.PrintStats()
}

func TestLRU(t *testing.T) {
	countCleaned := 0
	countAdded := 0

	c := lru.NewSafeLRUCache(cacheSize)
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
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i ++ {
				c.Get(lines[i])
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Fatalf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	c.PrintStats()
}

func TestRRC(t *testing.T) {
	countCleaned := 0
	countAdded := 0

	c := rrc.NewSafeRRCache(cacheSize)
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
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i ++ {
				c.Get(lines[i])
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	c.PrintStats()
}