package cache

// NewStore creates a new Store.
func NewStore[T comparable](keyFunc KeyFunc[T]) Store[T] {
	return &cache[any, T]{
		store:   NewThreadSafeStore(Indexers[any]{}, Indexes[any, T]{}),
		keyFunc: keyFunc,
	}
}

// NewIndexer creates a new IndexedStore.
func NewIndexer[K, T comparable](keyFunc KeyFunc[T]) IndexedStore[K, T] {
	return &cache[K, T]{
		store:   NewThreadSafeStore(Indexers[K]{}, Indexes[K, T]{}),
		keyFunc: keyFunc,
	}
}

// cache implements Store and IndexedStore.
type cache[K, T comparable] struct {
	store ThreadSafeStore[K, T]
	// keyFunc is used to make the key for objects stored in and retrieved from items
	keyFunc KeyFunc[T]
}

var _ Store[any] = &cache[any, any]{}
var _ IndexedStore[any, any] = &cache[any, any]{}

// Add inserts an item into the cache.
func (c *cache[K, T]) Add(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	c.store.Add(key, obj)
	return nil
}

// Update sets an item in the cache to its updated state.
func (c *cache[K, T]) Update(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	c.store.Update(key, obj)
	return nil
}

// Delete removes an item from the cache.
func (c *cache[K, T]) Delete(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	c.store.Delete(key)
	return nil
}

// List returns a list of all the items.
func (c *cache[K, T]) List() []interface{} {
	return c.store.List()
}

// ListKeys returns a list of all the keys of the objects currently
// in the cache.
func (c *cache[K, T]) ListKeys() []T {
	return c.store.ListKeys()
}

// ListKeysByIndex returns the storage keys of the stored objects whose set of
// indexed values for the named index includes the given indexed value.
func (c *cache[K, T]) ListKeysByIndex(indexName string, indexedValue K) ([]T, error) {
	return c.store.IndexKeys(indexName, indexedValue, nil)
}

// ListByIndex returns the stored objects whose set of indexed values
// for the named index includes the given indexed value.
func (c *cache[K, T]) ListByIndex(indexName string, indexedValue K) ([]interface{}, error) {
	return c.store.ByIndex(indexName, indexedValue, nil)
}

// AddIndexer add new indexer.
func (c *cache[K, T]) AddIndexer(indexName string, indexFunc IndexFunc[K]) error {
	return c.store.AddIndexer(indexName, indexFunc)
}

// AddIndexers adds more indexers to this store.
func (c *cache[K, T]) AddIndexers(newIndexers Indexers[K]) error {
	return c.store.AddIndexers(newIndexers)
}

// Get returns the requested item。
func (c *cache[K, T]) Get(obj interface{}) (item interface{}, exists bool, err error) {
	key, err := c.keyFunc(obj)
	if err != nil {
		return nil, false, KeyError{obj, err}
	}
	return c.GetByKey(key)
}

// GetByKey returns the requested item。
func (c *cache[K, T]) GetByKey(key T) (interface{}, bool, error) {
	item, exists := c.store.Get(key)
	return item, exists, nil
}

// Replace will delete the contents of 'c', using instead the given list.
func (c *cache[K, T]) Replace(list []interface{}) error {
	items := make(map[T]interface{}, len(list))
	for _, item := range list {
		key, err := c.keyFunc(item)
		if err != nil {
			return KeyError{item, err}
		}
		items[key] = item
	}
	c.store.Replace(items)
	return nil
}

// Size returns count of object in the cache.
func (c *cache[K, T]) Size() int {
	return c.store.Size()
}
