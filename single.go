package cache

import (
	"fmt"
	"sync"
)

// Single is a cache that saves a value at most.
type Single[K comparable, V any] struct {
	sync.RWMutex
	Stat
	source Source[K, V]

	key   K
	value V
}

// NewSingle returns a new Single cache.
// If a cache miss occurs, this tries to get a value from `source`.
func NewSingle[K comparable, V any](source Source[K, V]) (*Single[K, V], error) {
	if source == nil {
		return nil, ErrNoSource
	}
	return &Single[K, V]{
		source: source,
	}, nil
}

func (s *Single[K, V]) Get(key K) (V, error) {
	s.Lock()
	defer s.Unlock()

	if key == s.key {
		s.hit++
		return s.value, nil
	}
	s.miss++

	value, err := s.source(key)
	if err != nil {
		var v V
		return v, err
	}
	s.size = 1
	s.key = key
	s.value = value
	return value, nil
}

func (s *Single[K, V]) String() string {
	return fmt.Sprintf("%v %v", s.key, s.value)
}
