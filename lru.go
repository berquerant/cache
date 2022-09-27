package cache

import (
	"fmt"
	"strings"
	"sync"
)

// LRU implements a Least-Recentry-Uses cache.
type LRU[K comparable, V any] struct {
	sync.RWMutex
	Stat
	db      map[K]*lruEntry[K, V]
	source  Source[K, V]
	maxSize int

	// linked list
	head *lruEntry[K, V] // latest entry, prev is nil
	tail *lruEntry[K, V] // oldest entry, next is nil
}

// NewLRU returns a new LRU cache.
// It can save `size` values at most.
// If a cache miss occurs, this tries to get a value from `source`.
func NewLRU[K comparable, V any](size int, source Source[K, V]) (*LRU[K, V], error) {
	if source == nil {
		return nil, ErrNoSource
	}
	if size < 2 {
		return nil, fmt.Errorf("%w LRU size must be greater than 1", ErrInvalidSize)
	}
	return &LRU[K, V]{
		db:      make(map[K]*lruEntry[K, V]),
		source:  source,
		maxSize: size,
	}, nil
}

type lruEntry[K comparable, V any] struct {
	key   K
	value V
	next  *lruEntry[K, V]
	prev  *lruEntry[K, V]
}

// Remove tail element and delete it from db.
func (lru *LRU[K, V]) removeTailWithoutLock() {
	switch {
	case lru.head == nil: // len == 0
		return
	case lru.head == lru.tail: // len == 1
		lru.size--
		delete(lru.db, lru.tail.key)
		lru.head = nil
		lru.tail = nil
	default:
		lru.size--
		delete(lru.db, lru.tail.key)
		lru.tail = lru.tail.prev
		lru.tail.next = nil
	}
}

// Insert a new element to head.
func (lru *LRU[K, V]) insertToHeadWithoutLock(key K, value V) {
	lru.size++
	entry := &lruEntry[K, V]{
		key:   key,
		value: value,
	}

	lru.db[key] = entry
	switch {
	case lru.head == nil: // len == 0
		lru.head = entry
		lru.tail = entry
	default:
		oldHead := lru.head
		entry.next = oldHead
		oldHead.prev = entry
		lru.head = entry
	}
}

func (lru *LRU[K, V]) moveToHeadWithoutLock(entry *lruEntry[K, V]) {
	switch {
	case entry == nil || lru.head == nil || lru.head == lru.tail || lru.head == entry:
		// len == 0, 1 or already head
	case lru.tail == entry:
		// move tail to head
		lru.removeTailWithoutLock()
		lru.insertToHeadWithoutLock(entry.key, entry.value)
	default:
		// prev => entry => next
		// into
		// prev => next
		prev := entry.prev
		next := entry.next
		prev.next = next
		next.prev = prev
		lru.size--
		delete(lru.db, entry.key)
		lru.insertToHeadWithoutLock(entry.key, entry.value)
	}
}

func (lru *LRU[K, V]) Get(key K) (V, error) {
	lru.Lock()
	defer lru.Unlock()

	if entry, found := lru.db[key]; found {
		lru.hit++
		if lru.head == lru.tail { // len == 1
			return entry.value, nil
		}
		// move to head because entry is accssed right now
		lru.moveToHeadWithoutLock(entry)
		return entry.value, nil
	}
	lru.miss++

	value, err := lru.source(key)
	if err != nil {
		var v V
		return v, err
	}

	lru.insertToHeadWithoutLock(key, value)
	if lru.size > lru.maxSize {
		lru.removeTailWithoutLock()
	}
	return value, nil
}

func (lru *LRU[K, V]) String() string {
	lru.RLock()
	defer lru.RUnlock()

	var ss []string
	for e := lru.head; e != nil; e = e.next {
		ss = append(ss, e.String())
	}
	return fmt.Sprintf("tail %v list %s db %v", lru.tail, strings.Join(ss, " "), lru.db)
}

func (e *lruEntry[K, V]) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%v => %v prev ", e.key, e.value))
	if e.prev != nil {
		b.WriteString(fmt.Sprint(e.prev.key))
	} else {
		b.WriteString("nil")
	}
	b.WriteString(" next ")
	if e.next != nil {
		b.WriteString(fmt.Sprint(e.next.key))
	} else {
		b.WriteString("nil")
	}
	return fmt.Sprintf("[%s]", b.String())
}
