package test

import "testing"
//import "go-cache/arc"
import "go-cache/lru"
import "go-cache/base"
import "strings"
import "io/ioutil"
import "sync"
//import "time"

/*
func TestARCFetch(t *testing.T) {
	c := arc.NewSafeArcCache(cacheSize)
	data, err := ioutil.ReadFile("list.txt")
	if err != nil {
		t.Errorf("err: %s\n", err)
	}
	str := string(data)
	lines := strings.Split(str, "\n")

	c.SetFetchFunc(fetch)

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

	for key, obj := range(c.GetAllObjects()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	c.PrintStats()
}
*/
func TestLRUFetch(t *testing.T) {
	c := lru.NewLRUCache(cacheSize)
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
			for i := 0; i < countAccess; i ++ {
				_, err := c.Get(lines[i])
				if err != nil {
					c.Set(lines[i], &StringObject{s:lines[i]})
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.Check()

	for key, obj := range(c.Collect()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
}

func TestRandomFetch(t *testing.T) {
	c := base.NewRRCache(cacheSize)
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
			for i := 0; i < countAccess; i ++ {
				_, err := c.Get(lines[i])
				if err != nil {
					c.Set(lines[i], &StringObject{s:lines[i]})
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	c.Check()

	for key, obj := range(c.Collect()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
}