package test

import "testing"
import "strconv"
import "go-cache"
import "go-cache/arc"
import "math/rand"

func TestGet(t *testing.T) {
	cacheSize := 100
	countCleaned := 0
	countAdded := 0

	arc := arc.NewArcCache(cacheSize)

	arc.SetCleanFunc(func (obj cache.CacheObject) error {
		countCleaned += 1
		return nil
	})
	arc.SetFetchFunc(func (key string) (cache.CacheObject, error) {
		countAdded += 1
		return key, nil
	})

	for i := 0; i < 2000; i ++ {
		arc.Get("key"+strconv.Itoa(i))
	}

	for i := 0; i < 200; i ++ {
		j := rand.Intn(20)
		arc.Get("key"+strconv.Itoa(j))
	}

	if countCleaned + cacheSize != countAdded {
		t.Errorf("numbers of data items dont match")
	}
	arc.Checkkeys()
}