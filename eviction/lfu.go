package eviction

import (
	"container/heap"
	"sync"
)

// LFU implements the Least Frequently Used eviction policy.
type LFU[T comparable] struct {
	mu       sync.Mutex
	capacity int
	cache    map[T]*lfuEntry[T]
	freqHeap *lfuHeap[T]
}

type lfuEntry[T comparable] struct {
	key       T
	frequency int
	index     int
}

type lfuHeap[T comparable] []*lfuEntry[T]

// NewLFU creates a new LFU cache with the given capacity.
func NewLFU[T comparable](capacity int) Policy[T] {
	return &LFU[T]{
		capacity: capacity,
		cache:    make(map[T]*lfuEntry[T]),
		freqHeap: &lfuHeap[T]{},
	}
}

// Put adds a key to the cache. If the cache is full, it evicts the least frequently used key.
func (l *LFU[T]) Put(key T) (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var evictedKey T
	var evicted bool

	if entry, ok := l.cache[key]; ok {
		entry.frequency++
		heap.Fix(l.freqHeap, entry.index)
		return evictedKey, false
	}
	if len(l.cache) >= l.capacity {
		evictedKey, evicted = l.evict()
	}
	entry := &lfuEntry[T]{key: key, frequency: 1}
	heap.Push(l.freqHeap, entry)
	l.cache[key] = entry
	return evictedKey, evicted
}

// Delete removes a key from the cache.
func (l *LFU[T]) Delete(key T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if entry, ok := l.cache[key]; ok {
		heap.Remove(l.freqHeap, entry.index)
		delete(l.cache, key)
	}
}

// Reset clears all keys from the cache.
func (l *LFU[T]) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cache = make(map[T]*lfuEntry[T])
	l.freqHeap = &lfuHeap[T]{}
}

// Size returns the current number of keys in the cache.
func (l *LFU[T]) Size() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return len(l.cache)
}

// Evict removes the least frequently used key from the cache.
func (l *LFU[T]) Evict() (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.evict()
}

// evict is an internal method that removes the least frequently used key from the cache.
func (l *LFU[T]) evict() (T, bool) {
	if len(*l.freqHeap) == 0 {
		var zero T
		return zero, false
	}
	entry := heap.Pop(l.freqHeap).(*lfuEntry[T])
	delete(l.cache, entry.key)
	return entry.key, true
}

func (h lfuHeap[T]) Len() int           { return len(h) }
func (h lfuHeap[T]) Less(i, j int) bool { return h[i].frequency < h[j].frequency }
func (h lfuHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *lfuHeap[T]) Push(x interface{}) {
	entry := x.(*lfuEntry[T])
	entry.index = len(*h)
	*h = append(*h, entry)
}
func (h *lfuHeap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	entry := old[n-1]
	*h = old[0 : n-1]
	return entry
}
