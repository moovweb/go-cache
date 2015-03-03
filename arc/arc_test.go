package arc

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

// Taken from MPS test code that fails
// https://github.com/moovweb/manhattan/blob/773b5c6820ea65c43846649d9dfa01016e08cae8/targets/mps/main_test.go#L237
// It uncovers the issue with go-cache bug PI-202
// TODO: Redo the test, and add proper assertions. Here we just wait until the code panics
// TODO: remove refs to private repos
func TestReproEvictionPanic(t *testing.T) {
	var err error
	cacheSizeUnit := 1024 * 1024
	c := NewSafeARCache(cacheSizeUnit)
	s1 := string(make([]byte, cacheSizeUnit))
	obj1 := &StringObject{s: s1}
	s2 := string(make([]byte, cacheSizeUnit))
	obj2 := &StringObject{s: s2}
	s3 := string(make([]byte, cacheSizeUnit+1))
	too_big_obj := &StringObject{s: s3}

	err = c.Set("key1", too_big_obj)
	_, err = c.Get("key1")
	err = c.Set("key1", obj1)
	err = c.Set("key2", obj2)
	if err != nil {
		println(err.Error())
	}

}

// Turns out calling Get() again is having an effect on cache
// (The cache then reblances blocks or something like that)
// TODO: remove this test after you fix the bug and improve previous test case
func TestEvictionPanicHackyFix(t *testing.T) {
	var err error
	cacheSizeUnit := 1024 * 1024
	c := NewSafeARCache(cacheSizeUnit)
	s1 := string(make([]byte, cacheSizeUnit))
	obj1 := &StringObject{s: s1}
	s2 := string(make([]byte, cacheSizeUnit))
	obj2 := &StringObject{s: s2}
	s3 := string(make([]byte, cacheSizeUnit+1))
	too_big_obj := &StringObject{s: s3}

	println("before 1st Set")
	err = c.Set("key1", too_big_obj)
	println("before 1st Get")
	_, err = c.Get("key1")
	println("before 2nd Set")
	err = c.Set("key1", obj1)
	println("before 3rd Set")
	err = c.Set("key2", obj2)
	if err != nil {
		println(err.Error())
	}
}

type TestCacheObject struct {
}
