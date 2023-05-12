package ty

import (
	"encoding/json"
	"errors"
	"time"
)

type MI map[string]interface{}
type MS map[string]string

const Format = time.RFC3339

func (mi MI) GetString(key string) string {
	if v, b := mi[key]; b {
		return v.(string)
	}
	return ""
}

func (mi MI) GetBool(key string) bool {
	if v, b := mi[key]; b {
		return v.(bool)
	}
	return false
}

func MergeM[T interface{}](parent map[string]T, child map[string]T) map[string]T {
	for k, v := range child {
		parent[k] = v
	}

	return parent
}

type Opt[T interface{}] struct {
	Value T // inner value
	Set   bool
	Valid bool
}

func OptWrap[T interface{}](value T) Opt[T] {
	return Opt[T]{
		Value: value,
		Set:   true,
		Valid: true,
	}
}

func (i *Opt[T]) Merge(or *Opt[T]) {
	if or.Set {
		i.Value = or.Value
		i.Set = or.Set
		i.Valid = or.Valid
	}
}

func (i *Opt[T]) S(v T) {
	i.Value = v
	i.Set = true
	i.Valid = true
}

func (i *Opt[T]) N() {
	i.Valid = false
}

func (i *Opt[T]) U() {
	i.Set = false
	i.Valid = false
}

func (i *Opt[T]) UnmarshalJSON(data []byte) error {
	i.Set = true

	if string(data) == "null" {
		i.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &i.Value); err != nil {
		return err
	}

	i.Valid = true

	return nil
}

func (i *Opt[T]) MarshalJSON() ([]byte, error) {
	if !i.Set {
		return []byte("null"), nil
	}
	if !i.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(i.Value)
}

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
