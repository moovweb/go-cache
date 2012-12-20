package base

type CacheDirectoryBlock interface {
	GetKey() string
	SetKey(string)
	GetEntry() *CacheEntry
	SetEntry(*CacheEntry)
	IsEntryCached() bool
	IsEntryNil() bool
}

type BaseCdb struct {
	entry *CacheEntry
	key string
}

func (cdb *BaseCdb) GetKey() string {
	return cdb.key
}

func (cdb *BaseCdb) SetKey(key string) {
	cdb.key = key
}

func (cdb *BaseCdb) IsEntryCached() bool {
	return cdb.entry != nil
}

func (cdb *BaseCdb) SetEntry(entry *CacheEntry) {
	cdb.entry = entry
}

func (cdb *BaseCdb) GetEntry() *CacheEntry {
	return cdb.entry
}

func (cdb *BaseCdb) IsEntryNil() bool {
	return cdb.entry == nil
}