package cache

import (
	"fmt"
	"strings"
	"sync"
)

// FIFO implements a FIFO cache.
type FIFO[K comparable, V any] struct {
	sync.RWMutex
	Stat
	source Source[K, V]
	db     map[K]*fifoCell[K, V]

	// ring buffer
	cells []*fifoCell[K, V]
	head  int // index of a cell to be written, the oldest cell
}

// NewFIFO returns a new FIFO cache.
// It can save `size` values at most.
// If a cache miss occurs, this tries to get a value from `source`.
func NewFIFO[K comparable, V any](size int, source Source[K, V]) (*FIFO[K, V], error) {
	if source == nil {
		return nil, ErrNoSource
	}
	if size < 2 {
		return nil, fmt.Errorf("%w FIFO size must be greater than 1", ErrInvalidSize)
	}
	return &FIFO[K, V]{
		source: source,
		cells:  make([]*fifoCell[K, V], size),
		db:     make(map[K]*fifoCell[K, V]),
	}, nil
}

func (f *FIFO[K, V]) Get(key K) (V, error) {
	f.Lock()
	defer f.Unlock()

	if cell, found := f.db[key]; found {
		f.hit++
		return cell.value, nil
	}
	f.miss++

	value, err := f.source(key)
	if err != nil {
		var v V
		return v, err
	}

	if oldValue := f.cells[f.getHead()]; oldValue != nil {
		delete(f.db, oldValue.key)
	} else {
		f.size++
	}

	cell := &fifoCell[K, V]{
		key:   key,
		value: value,
	}
	f.db[key] = cell
	f.cells[f.getHead()] = cell
	f.incrHead()
	return value, nil
}

func (f *FIFO[K, V]) getHead() int { return f.index(f.head) }
func (f *FIFO[K, V]) incrHead()    { f.head = f.index(f.head + 1) }
func (f *FIFO[K, V]) index(i int) int {
	for i < 0 {
		i += len(f.cells)
	}
	return i % len(f.cells)
}

func (f *FIFO[K, V]) String() string {
	f.RLock()
	defer f.RUnlock()

	var ss []string
	for _, x := range f.cells {
		if x != nil {
			ss = append(ss, fmt.Sprintf("[%v]", x))
		}
	}
	return fmt.Sprintf("list %s db %v", strings.Join(ss, " "), f.db)
}

type fifoCell[K comparable, V any] struct {
	key   K
	value V
}

func (c *fifoCell[K, V]) String() string {
	return fmt.Sprintf("%v => %v", c.key, c.value)
}
