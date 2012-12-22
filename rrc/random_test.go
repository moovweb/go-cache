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
	return 1
}

func TestGet(t *testing.T) {
	cacheSize := 100
	countAdded := 0
	countCleaned := 0
	countAccess := 2000

	c := NewSafeRRCache(cacheSize)

	c.SetCleanFunc(func (obj cache.CacheObject) error {
		countCleaned += 1
		println("replacing", obj.(*StringObject).s)
		return nil
	})
	rand.Seed(time.Now().Unix())

	for i := 0; i < countAccess; i ++ {
		j := rand.Intn(cacheSize*2)
		key := "key"+strconv.Itoa(j)
		val, err := c.Get(key)
		if err == cache.CacheMiss {
			countAdded += 1
			c.Set(key, &StringObject{s: key})
		} else if val.(*StringObject).s != key {
			t.Errorf("key does not match the value")
		}
	}

	c.Check()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(c.Collect()) {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	c.PrintStats()
}