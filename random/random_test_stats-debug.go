// +build cache_debug

package rrc

func PrintStats(c *RRCache) {
	println("cache hit rate:", c.GetHitRate())
	println("cache avg access time:", c.GetAvgAccessTime())
}