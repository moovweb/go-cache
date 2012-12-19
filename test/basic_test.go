package test

import "testing"
import "strconv"
import "go-cache"
import "go-cache/arc"
import "math/rand"
import "time"

func TestGet(t *testing.T) {
	cacheSize := 100
	countCleaned := 0
	countAdded := 0
	countAccess := 2000

	arc := arc.NewArcCache(cacheSize)

	arc.SetCleanFunc(func (obj cache.CacheObject) error {
		countCleaned += 1
		return nil
	})
	arc.SetFetchFunc(func (key string) (cache.CacheObject, error) {
		countAdded += 1
		return key, nil
	})

	rand.Seed(time.Now().Unix())

	for i := 0; i < countAccess; i ++ {
		j := rand.Intn(cacheSize*2)
		val, err := arc.Get("key"+strconv.Itoa(j))
		if err != cache.CacheMiss && err != nil {
			t.Errorf("unexpected err:", err.Error())
		}
		if val == nil || val.(string) != "key"+strconv.Itoa(j) {
			t.Errorf("key does not match the value")
		}
	}

	countMiss := 0
	for i := 0; i < countAccess; i ++ {
		j := rand.Intn(cacheSize*2)
		val, err := arc.Get("key"+strconv.Itoa(j))
		if err == cache.CacheMiss {
			countMiss += 1
		}
		if val == nil || val.(string) != "key"+strconv.Itoa(j) {
			t.Errorf("key does not match the value")
		}
	}
	println("cache hit rate:", (100*(countAccess-countMiss))/countAccess)

	arc.CheckCache()

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match: %d != %d + %d\n", countAdded, countCleaned, cacheSize)
	}
	
	for key, obj := range(arc.GetAllObjects()) {
		if key != obj.(string) {
			t.Errorf("key does not match the cached value")
		}
	}
}