package eviction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	cache := NewLRU[int](2)

	// Test Put and Size
	evictedKey, evicted := cache.Put(1)
	assert.False(t, evicted)
	assert.Equal(t, 0, evictedKey)
	assert.Equal(t, 1, cache.Size())

	evictedKey, evicted = cache.Put(2)
	assert.False(t, evicted)
	assert.Equal(t, 0, evictedKey)
	assert.Equal(t, 2, cache.Size())

	// Test Put with eviction
	evictedKey, evicted = cache.Put(3)
	assert.True(t, evicted)
	assert.Equal(t, 1, evictedKey)
	assert.Equal(t, 2, cache.Size())

	// Test Delete
	cache.Delete(2)
	assert.Equal(t, 1, cache.Size())

	// Test Reset
	cache.Reset()
	assert.Equal(t, 0, cache.Size())

	// Test Evict
	cache.Put(1)
	cache.Put(2)
	key, ok := cache.Evict()
	assert.True(t, ok)
	assert.Equal(t, 1, key)
	assert.Equal(t, 1, cache.Size())
}

func TestLRUMultiEvictions(t *testing.T) {
	cache := NewLRU[int](3)

	// Fill the cache
	cache.Put(1)
	cache.Put(2)
	cache.Put(3)
	assert.Equal(t, 3, cache.Size())

	// Access some elements to change their order
	cache.Put(2)
	cache.Put(1)

	// Add another element to trigger eviction
	evictedKey, evicted := cache.Put(4)
	assert.True(t, evicted)
	assert.Equal(t, 3, evictedKey)
	assert.Equal(t, 3, cache.Size())
}

func TestLRUDeleteNonExistentKey(t *testing.T) {
	cache := NewLRU[int](10)

	// Delete non-existent key
	cache.Delete(1)
	assert.Equal(t, 0, cache.Size())

	// Add and then delete a key
	cache.Put(1)
	cache.Delete(1)
	assert.Equal(t, 0, cache.Size())
}
