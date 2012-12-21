package rrc

import "testing"
import "go-cache"
import "math/rand"
import "time"
import "strconv"

type StringObject struct {
	s string
}

func (o *StringObject) Size() int {
	return len(o.s)
}

func TestGet(t *testing.T) {
	cacheSize := 100
	countAdded := 0
	countCleaned := 0
	countAccess := 2000

	c := NewSafeRRCache(cacheSize)

	c.SetCleanFunc(func (obj cache.CacheObject) error {
		countCleaned += 1
		//println("replacing", obj.(string))
		return nil
	})
	c.SetFetchFunc(func (key string) (cache.CacheObject, error) {
		countAdded += 1
		return &StringObject{s:key}, nil
	})
	rand.Seed(time.Now().Unix())

	for i := 0; i < countAccess; i ++ {
		j := rand.Intn(cacheSize*2)
		val, err := c.Get("key"+strconv.Itoa(j))
		if err != cache.CacheMiss && err != nil {
			t.Errorf("unexpected err:", err.Error())
		}
		if val == nil || val.(*StringObject).s != "key"+strconv.Itoa(j) {
			t.Errorf("key does not match the value")
		}
	}

	for i := 0; i < countAccess; i ++ {
		j := rand.Intn(cacheSize*2)
		val, _ := c.Get("key"+strconv.Itoa(j))
		if val == nil || val.(*StringObject).s != "key"+strconv.Itoa(j) {
			t.Errorf("key does not match the value")
		}
	}

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