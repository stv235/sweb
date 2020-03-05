package database

type LocalCache struct {
	lookup func(id int64) interface{}
	items map[int64]interface{}
}

func (cache LocalCache) Lookup(id int64) interface{} {
	item, ok := cache.items[id]

	if !ok {
		item = cache.lookup(id)
		cache.items[id] = item
	}

	return item
}

func NewLocalCache(lookup func(id int64) interface{}) LocalCache {
	return LocalCache{
		lookup: lookup,
		items:  make(map[int64]interface{}),
	}
}
