package test

import "testing"
import "github.com/moovweb/go-cache/arc"
import "github.com/moovweb/go-cache/lru"
import "github.com/moovweb/go-cache/base"
import "strings"
import "io/ioutil"
import "sync"

//import "time"

func TestArcFetch(t *testing.T) {
	c := arc.NewSafeARCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	countAccess := len(lines)
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i++ {
				_, err := c.Get(lines[i])
				if err != nil {
					c.Set(lines[i], &StringObject{s: lines[i]})
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.Check()

	for key, obj := range c.Collect() {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
}
func TestLRUFetch(t *testing.T) {
	c := lru.NewSafeLRUCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	countAccess := len(lines)
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i++ {
				_, err := c.Get(lines[i])
				if err != nil {
					c.Set(lines[i], &StringObject{s: lines[i]})
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.Check()

	for key, obj := range c.Collect() {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
}

func TestRandomFetch(t *testing.T) {
	c := base.NewSafeRRCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	countAccess := len(lines)
	countAccess = 2000
	wg := &sync.WaitGroup{}
	for j := 0; j < concurrency; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < countAccess; i++ {
				_, err := c.Get(lines[i])
				if err != nil {
					c.Set(lines[i], &StringObject{s: lines[i]})
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.Check()

	for key, obj := range c.Collect() {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
}
