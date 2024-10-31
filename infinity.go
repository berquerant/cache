package cache

import "sync"

// Infinity is an unlimited cache.
type Infinity[K comparable, V any] struct {
	sync.Mutex
	Stat
	db     map[K]V
	source Source[K, V]
}

func NewInfinity[K comparable, V any](source Source[K, V]) *Infinity[K, V] {
	return &Infinity[K, V]{
		db:     make(map[K]V),
		source: source,
	}
}

func (f *Infinity[K, V]) Get(key K) (V, error) {
	f.Lock()
	defer f.Unlock()

	if v, ok := f.db[key]; ok {
		f.hit++
		return v, nil
	}
	f.miss++

	value, err := f.source(key)
	if err != nil {
		var v V
		return v, err
	}

	f.size++
	f.db[key] = value
	return value, nil
}
