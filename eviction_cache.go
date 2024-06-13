package cache

import (
	"fmt"
	"sync"

	"github.com/liuxinbot/cache/eviction"
)

// EvictionStore extends IndexedStore with eviction capabilities.
type EvictionStore[K, T comparable] interface {
	IndexedStore[K, T]

	Evict() error
}

// NewEvictionCache creates a new EvictionStore.
func NewEvictionCache[K comparable, T comparable](keyFunc KeyFunc[T], evictionPolicy eviction.Policy[T], indexers Indexers[K]) EvictionStore[K, T] {
	return &evictionCache[K, T]{
		store:          NewThreadSafeStore(indexers, make(Indexes[K, T])),
		keyFunc:        keyFunc,
		evictionPolicy: evictionPolicy,
	}
}

// cache implements IndexedStore and EvictionStore.
type evictionCache[K comparable, T comparable] struct {
	store          ThreadSafeStore[K, T]
	keyFunc        KeyFunc[T]
	evictionPolicy eviction.Policy[T]
	mu             sync.Mutex
}

// Add adds an object to the cache.
func (c *evictionCache[K, T]) Add(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Call Add on eviction policy
	evictedKey, evicted := c.evictionPolicy.Put(key)
	if evicted {
		// EvictionPolicy.Add returned true, indicating eviction occurred
		c.store.Delete(evictedKey) // Delete the eliminated key from store
	}

	// Add the new object to store
	c.store.Add(key, obj)
	return nil
}

// Update updates an object in the cache.
func (c *evictionCache[K, T]) Update(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store.Update(key, obj)
	c.evictionPolicy.Put(key)
	return nil
}

// Delete deletes an object from the cache.
func (c *evictionCache[K, T]) Delete(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictionPolicy.Delete(key)
	c.store.Delete(key)
	return nil
}

// List returns a list of all cached objects.
func (c *evictionCache[K, T]) List() []interface{} {
	return c.store.List()
}

// ListKeys returns a list of keys for all cached objects.
func (c *evictionCache[K, T]) ListKeys() []T {
	return c.store.ListKeys()
}

// ListKeysByIndex returns a list of keys based on the index name and indexed value.
func (c *evictionCache[K, T]) ListKeysByIndex(indexName string, indexedValue K) ([]T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys, err := c.store.IndexKeys(indexName, indexedValue, nil)
	if err != nil {
		return keys, err
	}
	for _, key := range keys {
		c.evictionPolicy.Put(key)
	}
	return keys, nil
}

// ListByIndex returns a list of objects based on the index name and indexed value.
func (c *evictionCache[K, T]) ListByIndex(indexName string, indexedValue K) ([]interface{}, error) {
	return c.store.ByIndex(indexName, indexedValue, nil)
}

// AddIndexer add new indexer.
func (c *evictionCache[K, T]) AddIndexer(indexName string, indexFunc IndexFunc[K]) error {
	return c.store.AddIndexer(indexName, indexFunc)
}

func (c *evictionCache[K, T]) AddIndexers(newIndexers Indexers[K]) error {
	return c.store.AddIndexers(newIndexers)
}

// Get retrieves an object from the cache based on the object.
func (c *evictionCache[K, T]) Get(obj interface{}) (interface{}, bool, error) {
	key, err := c.keyFunc(obj)
	if err != nil {
		return nil, false, KeyError{obj, err}
	}
	return c.GetByKey(key)
}

// GetByKey retrieves an object from the cache based on the key.
func (c *evictionCache[K, T]) GetByKey(key T) (interface{}, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, exists := c.store.Get(key)
	if exists {
		c.evictionPolicy.Put(key)
	}
	return item, exists, nil
}

// Replace replaces all objects in the cache.
func (c *evictionCache[K, T]) Replace(list []interface{}) error {
	items := make(map[T]interface{}, len(list))
	for _, item := range list {
		key, err := c.keyFunc(item)
		if err != nil {
			return KeyError{item, err}
		}
		items[key] = item
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// reset the eviction policy
	c.evictionPolicy.Reset()
	// Replace the store
	c.store.Replace(items)
	// Re-add items to eviction policy
	for key := range items {
		c.evictionPolicy.Put(key)
	}
	return nil
}

// Evict removes an object from the cache based on the cache eviction policy.
func (c *evictionCache[K, T]) Evict() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	key, ok := c.evictionPolicy.Evict()
	if !ok {
		return fmt.Errorf("no items to evict")
	}
	c.store.Delete(key)
	return nil
}

// Size returns count of object in the cache.
func (c *evictionCache[K, T]) Size() int {
	return c.store.Size()
}
