package uaxpl

import "sync/atomic"

var cache_ atomic.Pointer[Cache]

func init() {
	SetCache(&cache{})
}

func SetCache(cache Cache) {
	cache_.Store(&cache)
}

func GetCache() Cache {
	return *cache_.Load()
}
