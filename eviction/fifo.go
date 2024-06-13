package eviction

import (
	"container/list"
	"sync"
)

// FIFO implements the First In, First Out eviction policy.
type FIFO[T comparable] struct {
	mu       sync.Mutex
	capacity int
	cache    map[T]*list.Element
	list     *list.List
}

// NewFIFO creates a new FIFO cache with the given capacity.
func NewFIFO[T comparable](capacity int) Policy[T] {
	return &FIFO[T]{
		capacity: capacity,
		cache:    make(map[T]*list.Element),
		list:     list.New(),
	}
}

// Put adds a key to the cache. If the cache is full, it evicts the oldest key.
func (f *FIFO[T]) Put(key T) (T, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var evictedKey T
	var evicted bool

	if _, ok := f.cache[key]; ok {
		return evictedKey, false
	}
	if f.list.Len() >= f.capacity {
		evictedKey, evicted = f.evict()
	}
	elem := f.list.PushBack(&entry[T]{key})
	f.cache[key] = elem
	return evictedKey, evicted
}

// Delete removes a key from the cache.
func (f *FIFO[T]) Delete(key T) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if elem, ok := f.cache[key]; ok {
		f.list.Remove(elem)
		delete(f.cache, key)
	}
}

// Evict removes the oldest key from the cache.
func (f *FIFO[T]) Evict() (T, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.evict()
}

// Reset clears all keys from the cache.
func (f *FIFO[T]) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.cache = make(map[T]*list.Element)
	f.list.Init()
}

// Size returns the current number of keys in the cache.
func (f *FIFO[T]) Size() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return len(f.cache)
}

// evict is an internal method that removes the oldest key from the cache.
func (f *FIFO[T]) evict() (T, bool) {
	elem := f.list.Front()
	if elem == nil {
		var zero T
		return zero, false
	}
	f.list.Remove(elem)
	entry := elem.Value.(*entry[T])
	delete(f.cache, entry.key)
	return entry.key, true
}
