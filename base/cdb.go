package base

import . "github.com/moovweb/go-cache"

type CacheDirectoryBlock interface {
	GetKey() string
	SetKey(string)
	GetObject() CacheObject
	SetObject(CacheObject)
}

type BaseCdb struct {
	object CacheObject
	key    string
}

func (cdb *BaseCdb) GetKey() string {
	return cdb.key
}

func (cdb *BaseCdb) SetKey(key string) {
	cdb.key = key
}

func (cdb *BaseCdb) SetObject(object CacheObject) {
	cdb.object = object
}

func (cdb *BaseCdb) GetObject() CacheObject {
	return cdb.object
}

func NewBasicCdb() *BaseCdb {
	cdb := &BaseCdb{}
	return cdb
}
