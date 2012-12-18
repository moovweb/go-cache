package test

import "testing"
import "strconv"
import "go-cache"
import "go-cache/arc"

func clean(obj cache.CacheObject) error {
	println("clear:", obj.(string))
	return nil
}

func fetch(key string) (cache.CacheObject, error) {
	return key, nil
}

func TestGet(t *testing.T) {
	arc := arc.NewArcCache(10)
	arc.SetCleanFunc(clean)
	arc.SetFetchFunc(fetch)

	for i := 0; i < 20; i ++ {
		arc.Get("key"+strconv.Itoa(i))
	}
}