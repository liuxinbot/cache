package eviction

import (
	"container/list"
	"sync"
)

// lru implements the Least Recently Used eviction policy.
type lru[T comparable] struct {
	mu       sync.Mutex
	capacity int
	cache    map[T]*list.Element
	list     *list.List
}

type entry[T comparable] struct {
	key T
}

// NewLRU creates a new lru cache with the given capacity.
func NewLRU[T comparable](capacity int) Policy[T] {
	return &lru[T]{
		capacity: capacity,
		cache:    make(map[T]*list.Element),
		list:     list.New(),
	}
}

// Put adds a key to the cache. If the cache is full, it evicts the least recently used key.
func (l *lru[T]) Put(key T) (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var evictedKey T
	var evicted bool

	if elem, ok := l.cache[key]; ok {
		l.list.MoveToFront(elem)
		return evictedKey, false
	}
	if l.list.Len() >= l.capacity {
		evictedKey, evicted = l.evict()
	}
	elem := l.list.PushFront(&entry[T]{key})
	l.cache[key] = elem
	return evictedKey, evicted
}

// Delete removes a key from the cache.
func (l *lru[T]) Delete(key T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, ok := l.cache[key]; ok {
		l.list.Remove(elem)
		delete(l.cache, key)
	}
}

// Reset clears all keys from the cache.
func (l *lru[T]) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cache = make(map[T]*list.Element)
	l.list.Init()
}

// Size returns the current number of keys in the cache.
func (l *lru[T]) Size() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return len(l.cache)
}

// Evict removes the least recently used key from the cache.
func (l *lru[T]) Evict() (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.evict()
}

// evict is an internal method that removes the least recently used key from the cache.
func (l *lru[T]) evict() (T, bool) {
	elem := l.list.Back()
	if elem == nil {
		var zero T
		return zero, false
	}
	l.list.Remove(elem)
	entry := elem.Value.(*entry[T])
	delete(l.cache, entry.key)
	return entry.key, true
}
