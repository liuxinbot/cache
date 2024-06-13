package cache

import (
	"sync"
)

// ThreadSafeStore defines an interface for a thread-safe store with indexing capabilities.
type ThreadSafeStore[K, T comparable] interface {
	// Add an object to the store.
	Add(key T, obj interface{})

	// Update an object in the store.
	Update(key T, obj interface{})

	// Delete an object from the store.
	Delete(key T)

	// Get retrieve an object from the store.
	Get(key T) (item interface{}, exists bool)

	// List all objects in the store.
	List() []interface{}

	// ListKeys List all keys in the store.
	ListKeys() []T

	// Replace all objects in the store.
	Replace(items map[T]interface{})

	// Size get count of elements in the store.
	Size() int

	// Index retrieve objects by index.
	Index(indexName string, obj interface{}, lessFunc func(lhs T, rhs T) bool) ([]interface{}, error)

	// IndexKeys retrieve keys by index.
	IndexKeys(indexName string, indexedValue K, lessFunc func(lhs T, rhs T) bool) ([]T, error)

	// ByIndex retrieve objects by indexed value.
	ByIndex(indexName string, indexedValue K, lessFunc func(lhs, rhs T) bool) ([]interface{}, error)

	// AddIndexer add new indexer.
	AddIndexer(indexName string, indexFunc IndexFunc[K]) error

	// AddIndexers add new indexers.
	AddIndexers(newIndexers Indexers[K]) error
}

// threadSafeMap implements the ThreadSafeStore interface.
type threadSafeMap[K, T comparable] struct {
	mu    sync.RWMutex
	items map[T]interface{}
	index *storeIndex[K, T]
}

// NewThreadSafeStore creates a new instance of ThreadSafeStore.
func NewThreadSafeStore[K, T comparable](indexers Indexers[K], indices Indexes[K, T]) ThreadSafeStore[K, T] {
	return &threadSafeMap[K, T]{
		items: make(map[T]interface{}),
		index: &storeIndex[K, T]{
			indexers: indexers,
			indices:  indices,
		},
	}
}

// Add adds an object to the store.
func (tsm *threadSafeMap[K, T]) Add(key T, obj interface{}) {
	tsm.Update(key, obj)
}

// Update updates an object in the store.
func (tsm *threadSafeMap[K, T]) Update(key T, obj interface{}) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	oldObject := tsm.items[key]
	tsm.items[key] = obj
	tsm.index.updateIndices(oldObject, obj, key)
}

// Delete deletes an object from the store.
func (tsm *threadSafeMap[K, T]) Delete(key T) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	if obj, exists := tsm.items[key]; exists {
		tsm.index.updateIndices(obj, nil, key)
		delete(tsm.items, key)
	}
}

// Get retrieves an object from the store.
func (tsm *threadSafeMap[K, T]) Get(key T) (item interface{}, exists bool) {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()
	item, exists = tsm.items[key]
	return item, exists
}

// List lists all objects in the store.
func (tsm *threadSafeMap[K, T]) List() []interface{} {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()
	list := make([]interface{}, 0, len(tsm.items))
	for _, item := range tsm.items {
		list = append(list, item)
	}
	return list
}

// ListKeys lists all keys in the store.
func (tsm *threadSafeMap[K, T]) ListKeys() []T {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()
	list := make([]T, 0, len(tsm.items))
	for key := range tsm.items {
		list = append(list, key)
	}
	return list
}

// Replace replaces all objects in the store.
func (tsm *threadSafeMap[K, T]) Replace(items map[T]interface{}) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	tsm.items = items

	// Rebuild any index
	tsm.index.reset()
	for key, item := range tsm.items {
		tsm.index.updateIndices(nil, item, key)
	}
}

// Index retrieves objects by index.
func (tsm *threadSafeMap[K, T]) Index(indexName string, obj interface{}, lessFunc func(lhs, rhs T) bool) ([]interface{}, error) {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()

	keySet, err := tsm.index.getKeysFromIndex(indexName, obj)
	if err != nil {
		return nil, err
	}

	var keys []T
	if lessFunc == nil {
		keys = keySet.UnsortedList()
	} else {
		keys = keySet.List(lessFunc)
	}

	list := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		list = append(list, tsm.items[key])
	}
	return list, nil
}

// ByIndex retrieves objects by indexed value.
func (tsm *threadSafeMap[K, T]) ByIndex(indexName string, indexedValue K, lessFunc func(lhs, rhs T) bool) ([]interface{}, error) {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()

	keys, err := tsm.IndexKeys(indexName, indexedValue, lessFunc)
	if err != nil {
		return nil, err
	}

	list := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		list = append(list, tsm.items[key])
	}

	return list, nil
}

// IndexKeys retrieves keys by index.
func (tsm *threadSafeMap[K, T]) IndexKeys(indexName string, indexedValue K, lessFunc func(lhs, rhs T) bool) ([]T, error) {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()

	keySet, err := tsm.index.getKeysByIndex(indexName, indexedValue)
	if err != nil {
		return nil, err
	}

	if lessFunc == nil {
		return keySet.UnsortedList(), nil
	}

	return keySet.List(lessFunc), nil
}

// AddIndexers adds new indexers to the store.
func (tsm *threadSafeMap[K, T]) AddIndexers(newIndexers Indexers[K]) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	if err := tsm.index.addIndexers(newIndexers); err != nil {
		return err
	}

	// If there are already items, reindex them
	for key, item := range tsm.items {
		for name := range newIndexers {
			tsm.index.updateSingleIndex(name, nil, item, key)
		}
	}

	return nil
}

// AddIndexer adds new indexer to the store.
func (tsm *threadSafeMap[K, T]) AddIndexer(indexName string, indexFunc IndexFunc[K]) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	if err := tsm.index.addIndexer(indexName, indexFunc); err != nil {
		return err
	}

	// If there are already items, reindex them
	for key, item := range tsm.items {
		tsm.index.updateSingleIndex(indexName, nil, item, key)
	}

	return nil
}

// Size get count of elements in the store.
func (tsm *threadSafeMap[K, T]) Size() int {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	return len(tsm.items)
}
