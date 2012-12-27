package base

func NewRRCache(size int) * BaseCache {
	cdbm := NewBasicCdbm()
	c := NewBaseCache(size, cdbm)

	return c
}

func NewSafeRRCache(size int) * BaseCache {
	cdbm := NewBasicCdbm()
	c := NewSafeBaseCache(size, cdbm)

	return c
}