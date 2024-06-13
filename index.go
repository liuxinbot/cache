package cache

import (
	"fmt"

	"github.com/liuxinbot/cache/sets"
)

// IndexedStore extends Store with indexing capabilities.
type IndexedStore[K, T comparable] interface {
	Store[T]

	// ListKeysByIndex returns storage keys of objects whose indexed values for the specified index include the given indexed value.
	ListKeysByIndex(indexName string, indexedValue K) ([]T, error)

	// ListByIndex returns objects whose indexed values for the specified index include the given indexed value.
	ListByIndex(indexName string, indexedValue K) ([]interface{}, error)

	// AddIndexer add new indexer.
	AddIndexer(indexName string, indexFunc IndexFunc[K]) error

	// AddIndexers adds more indexers to this store.
	AddIndexers(newIndexers Indexers[K]) error
}

// IndexFunc is a function type that calculates a set of indexed values for an object.
type IndexFunc[K comparable] func(obj interface{}) ([]K, error)

// Index maps the indexed value to a set of keys in the store that match on that value.
type Index[K, T comparable] map[K]sets.Set[T]

// Indexers maps an index name to an IndexFunc.
type Indexers[K comparable] map[string]IndexFunc[K]

// Indexes maps an index name to an Index.
type Indexes[K, T comparable] map[string]Index[K, T]

// storeIndex implements the indexing functionality for a ThreadSafeStore.
type storeIndex[K, T comparable] struct {
	indexers Indexers[K]
	indices  Indexes[K, T]
}

// reset clears all indices.
func (si *storeIndex[K, T]) reset() {
	si.indices = Indexes[K, T]{}
}

// getKeysFromIndex retrieves the set of keys from the specified index that match the object.
func (si *storeIndex[K, T]) getKeysFromIndex(indexName string, obj interface{}) (sets.Set[T], error) {
	indexFunc, exists := si.indexers[indexName]
	if !exists {
		return nil, fmt.Errorf("index with name %s does not exist", indexName)
	}

	indexValues, err := indexFunc(obj)
	if err != nil {
		return nil, err
	}
	index := si.indices[indexName]

	var keySet sets.Set[T]
	if len(indexValues) == 1 {
		keySet = index[indexValues[0]]
	} else {
		keySet = sets.NewSet[T]()
		for _, indexValue := range indexValues {
			for key := range index[indexValue] {
				keySet.Insert(key)
			}
		}
	}
	return keySet, nil
}

// getKeysByIndex retrieves the set of keys from the specified index that match the indexed value.
func (si *storeIndex[K, T]) getKeysByIndex(indexName string, indexedValue K) (sets.Set[T], error) {
	_, exists := si.indexers[indexName]
	if !exists {
		return nil, fmt.Errorf("index with name %s does not exist", indexName)
	}
	index := si.indices[indexName]
	return index[indexedValue], nil
}

// addIndexer adds new indexer to the store.
func (si *storeIndex[K, T]) addIndexer(indexName string, indexFunc IndexFunc[K]) error {
	if _, exists := si.indexers[indexName]; exists {
		return fmt.Errorf("indexer conflict: %s", indexName)
	}
	si.indexers[indexName] = indexFunc
	return nil
}

// addIndexers adds new indexers to the store.
func (si *storeIndex[K, T]) addIndexers(newIndexers Indexers[K]) error {
	existingKeys := sets.KeySet[string](si.indexers)
	newKeys := sets.KeySet[string](newIndexers)

	if existingKeys.HasAny(newKeys.UnsortedList()...) {
		return fmt.Errorf("indexer conflict: %v", existingKeys.Intersection(newKeys))
	}

	for name, indexer := range newIndexers {
		si.indexers[name] = indexer
	}
	return nil
}

// updateIndices updates the object's location in the managed indexes:
// - For create, provide only the newObj
// - For update, provide both oldObj and newObj
// - For delete, provide only the oldObj
func (si *storeIndex[K, T]) updateIndices(oldObj, newObj interface{}, key T) {
	for name := range si.indexers {
		si.updateSingleIndex(name, oldObj, newObj, key)
	}
}

// updateSingleIndex updates a single index for the object.
func (si *storeIndex[K, T]) updateSingleIndex(name string, oldObj, newObj interface{}, key T) {
	var oldIndexValues, newIndexValues []K
	indexFunc, exists := si.indexers[name]
	if !exists {
		panic(fmt.Errorf("indexer %q does not exist", name))
	}

	if oldObj != nil {
		var err error
		oldIndexValues, err = indexFunc(oldObj)
		if err != nil {
			panic(fmt.Errorf("unable to calculate index entry for key %v on index %q: %v", key, name, err))
		}
	}

	if newObj != nil {
		var err error
		newIndexValues, err = indexFunc(newObj)
		if err != nil {
			panic(fmt.Errorf("unable to calculate index entry for key %v on index %q: %v", key, name, err))
		}
	}

	index := si.indices[name]
	if index == nil {
		index = Index[K, T]{}
		si.indices[name] = index
	}

	if len(newIndexValues) == 1 && len(oldIndexValues) == 1 && newIndexValues[0] == oldIndexValues[0] {
		return
	}

	for _, indexValue := range oldIndexValues {
		keySet := index[indexValue]
		if keySet == nil {
			return
		}
		keySet.Delete(key)
		if len(keySet) == 0 {
			delete(index, indexValue)
		}
	}
	for _, indexValue := range newIndexValues {
		keySet := index[indexValue]
		if keySet == nil {
			keySet = sets.NewSet[T]()
			index[indexValue] = keySet
		}
		keySet.Insert(key)
	}
}
