package ty

import "errors"

type Lazy[T interface{}] func() (*T, error)

func GetLazy[T interface{}](lazy func() (*T, error)) Lazy[T] {
	var cache *T
	return func() (*T, error) {
		if cache != nil {
			return cache, nil
		}
		cacheTmp, err := lazy()
		if err != nil {
			return cache, err
		}
		cache = cacheTmp
		return cache, nil
	}
}

type LazyMap[K string, V interface{}] map[K]Lazy[V]

func (lm LazyMap[K, V]) Get(key K) (*V, error) {
	val, ok := lm[key]
	if !ok {
		return nil, errors.New("not found " + string(key))
	}
	return val()
}
