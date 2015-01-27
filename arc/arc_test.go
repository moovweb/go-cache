package arc

import "testing"
import "github.com/moovweb/go-cache"
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
	cacheSize := 20
	countAdded := 0
	countCleaned := 0
	countAccess := 10000

	c := NewARCache(cacheSize * 5)

	c.SetCleanFunc(func(obj cache.CacheObject) error {
		countCleaned += obj.Size()
		return nil
	})
	rand.Seed(time.Now().Unix())

	for i := 0; i < countAccess; i++ {
		j := rand.Intn(cacheSize * 2)
		key := "key" + strconv.Itoa(j)
		val, err := c.Get(key)

		if err == cache.CacheMiss {
			countAdded += len(key)
			c.Set(key, &StringObject{s: key})
		} else if val.(*StringObject).s != key {
			t.Errorf("key does not match the value, %s != %s", val.(*StringObject).s, key)
		}
	}

	c.Check()
	if countCleaned+c.GetUsage() != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, c.GetUsage())
	}

	for key, obj := range c.Collect() {
		if key != obj.(*StringObject).s {
			t.Errorf("key does not match the cached value")
		}
	}
	println("cache hit rate:", c.GetHitRate())
	println("cache usage rate:", c.GetUsageRate())
	c.Reset()
	if c.GetUsage() != 0 {
		t.Errorf("after reset, cache usage should be zero")
	}
}
