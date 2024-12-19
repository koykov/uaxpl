package uaxpl

var cache_ Cacher[any]

func init() {
	cache_ = &cacheNaive[cacheEntry]{}
}

func SetCache[T any](cache Cacher[T]) {
	cache_ = cache
}
