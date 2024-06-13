package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liuxinbot/cache/eviction"
)

// Dummy key function for test
func testIntKeyFunc(obj interface{}) (int, error) {
	return obj.(int), nil
}

func TestEvictionCacheFIFO(t *testing.T) {
	fifo := eviction.NewFIFO[int](2)
	store := NewEvictionCache(testIntKeyFunc, fifo, make(Indexers[int]))

	// Test Add and Size
	err := store.Add(1)
	assert.NoError(t, err)
	err = store.Add(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, store.Size())

	// Test Add with eviction
	err = store.Add(3)
	assert.NoError(t, err)
	assert.Equal(t, 2, store.Size())
	_, exists, _ := store.Get(1)
	assert.False(t, exists)
}

func TestEvictionCacheLRU(t *testing.T) {
	lru := eviction.NewLRU[int](2)
	store := NewEvictionCache(testIntKeyFunc, lru, make(Indexers[int]))

	// Test Add and Size
	err := store.Add(1)
	assert.NoError(t, err)
	err = store.Add(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, store.Size())

	// Test Add with eviction
	err = store.Add(3)
	assert.NoError(t, err)
	assert.Equal(t, 2, store.Size())
	_, exists, _ := store.Get(1)
	assert.False(t, exists)

	// Test LRU behavior
	_, _, err = store.Get(2) // Access to make 2 recently used
	assert.NoError(t, err)
	err = store.Add(4) // This should evict key 3
	assert.NoError(t, err)
	_, exists, _ = store.Get(3)
	assert.False(t, exists)
}

func TestEvictionCacheLFU(t *testing.T) {
	lfu := eviction.NewLFU[int](2)
	store := NewEvictionCache(testIntKeyFunc, lfu, make(Indexers[int]))

	// Test Add and Size
	err := store.Add(1)
	assert.NoError(t, err)
	err = store.Add(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, store.Size())

	// Test Add with eviction
	err = store.Add(3)
	assert.NoError(t, err)
	assert.Equal(t, 2, store.Size())
	_, exists, _ := store.Get(1)
	assert.False(t, exists)

	// Test LFU behavior
	_, _, err = store.Get(2) // Access to increase frequency of 2
	assert.NoError(t, err)
	_, _, err = store.Get(2) // Access to further increase frequency of 2
	assert.NoError(t, err)
	err = store.Add(4) // This should evict key 3, as 2 has higher frequency
	assert.NoError(t, err)
	_, exists, _ = store.Get(3)
	assert.False(t, exists)
	_, exists, _ = store.Get(2)
	assert.True(t, exists)
}
