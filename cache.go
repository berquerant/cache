// Package cache provides a in-memory cache.
package cache

import "errors"

// Cache is a in-memory cache.
type Cache[K comparable, V any] interface {
	Get(key K) (V, error)
}

// Source is a data source of a cache.
type Source[K comparable, V any] func(K) (V, error)

// A simple cache statistics report.
type Stat struct {
	hit  int
	miss int
	size int
}

func (s *Stat) Hit() int  { return s.hit }
func (s *Stat) Miss() int { return s.miss }
func (s *Stat) Size() int { return s.size }

var (
	// ErrInvalidSize is returned by a constructor of `Cache` when the max size is not proper.
	ErrInvalidSize = errors.New("InvalidSize")
	// ErrNoSource is returned by a constructor of `Cache` when the source function is nil.
	ErrNoSource = errors.New("NoSource")
)
