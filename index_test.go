package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liuxinbot/cache/sets"
)

// TestStoreIndexAddIndexer tests adding indexers to the storeIndex
func TestStoreIndexAddIndexer(t *testing.T) {
	indexers := Indexers[string]{}
	si := &storeIndex[string, string]{
		indexers: indexers,
		indices:  Indexes[string, string]{},
	}

	// Add an indexer
	err := si.addIndexer("name", func(obj interface{}) ([]string, error) {
		return []string{obj.(string)}, nil
	})
	assert.Nil(t, err)

	// Add the same indexer again and expect an error
	err = si.addIndexer("name", func(obj interface{}) ([]string, error) {
		return []string{obj.(string)}, nil
	})
	assert.NotNil(t, err)
	assert.Equal(t, "indexer conflict: name", err.Error())
}

// TestStoreIndexGetKeysFromIndex tests retrieving keys by index
func TestStoreIndexGetKeysFromIndex(t *testing.T) {
	indexers := Indexers[string]{
		"name": func(obj interface{}) ([]string, error) {
			return []string{obj.(string)}, nil
		},
	}
	si := &storeIndex[string, string]{
		indexers: indexers,
		indices:  Indexes[string, string]{},
	}

	// Add objects
	si.updateIndices(nil, "obj1", "key1")
	si.updateIndices(nil, "obj2", "key2")

	// Retrieve keys from index
	keys, err := si.getKeysFromIndex("name", "obj1")
	assert.Nil(t, err)
	assert.Equal(t, sets.NewSet("key1"), keys)

	keys, err = si.getKeysFromIndex("name", "obj2")
	assert.Nil(t, err)
	assert.Equal(t, sets.NewSet("key2"), keys)
}

// TestStoreIndexGetKeysByIndex tests retrieving keys by indexed value
func TestStoreIndexGetKeysByIndex(t *testing.T) {
	indexers := Indexers[string]{
		"name": func(obj interface{}) ([]string, error) {
			return []string{obj.(string)}, nil
		},
	}
	si := &storeIndex[string, string]{
		indexers: indexers,
		indices:  Indexes[string, string]{},
	}

	// Add objects
	si.updateIndices(nil, "obj1", "key1")
	si.updateIndices(nil, "obj2", "key2")

	// Retrieve keys by index value
	keys, err := si.getKeysByIndex("name", "obj1")
	assert.Nil(t, err)
	assert.Equal(t, sets.NewSet("key1"), keys)

	keys, err = si.getKeysByIndex("name", "obj2")
	assert.Nil(t, err)
	assert.Equal(t, sets.NewSet("key2"), keys)
}

// TestStoreIndexUpdateIndices tests updating indices
func TestStoreIndexUpdateIndices(t *testing.T) {
	indexers := Indexers[string]{
		"name": func(obj interface{}) ([]string, error) {
			return []string{obj.(string)}, nil
		},
	}
	si := &storeIndex[string, string]{
		indexers: indexers,
		indices:  Indexes[string, string]{},
	}

	// Add objects
	si.updateIndices(nil, "obj1", "key1")
	si.updateIndices(nil, "obj2", "key2")

	// Update object
	si.updateIndices("obj1", "obj1_updated", "key1")

	// Verify updated indices
	keys, err := si.getKeysFromIndex("name", "obj1_updated")
	assert.Nil(t, err)
	assert.Equal(t, sets.NewSet("key1"), keys)

	keys, err = si.getKeysFromIndex("name", "obj1")
	assert.Nil(t, err)
	assert.Nil(t, keys)
}

// TestStoreIndexDeleteIndices tests deleting indices
func TestStoreIndexDeleteIndices(t *testing.T) {
	indexers := Indexers[string]{
		"name": func(obj interface{}) ([]string, error) {
			return []string{obj.(string)}, nil
		},
	}
	si := &storeIndex[string, string]{
		indexers: indexers,
		indices:  Indexes[string, string]{},
	}

	// Add objects
	si.updateIndices(nil, "obj1", "key1")
	si.updateIndices(nil, "obj2", "key2")

	// Delete object
	si.updateIndices("obj1", nil, "key1")

	// Verify deleted indices
	keys, err := si.getKeysFromIndex("name", "obj1")
	assert.Nil(t, err)
	assert.Nil(t, keys)
}

// TestStoreIndexAddIndexers tests adding multiple indexers
func TestStoreIndexAddIndexers(t *testing.T) {
	indexers := Indexers[string]{}
	si := &storeIndex[string, string]{
		indexers: indexers,
		indices:  Indexes[string, string]{},
	}

	newIndexers := Indexers[string]{
		"name": func(obj interface{}) ([]string, error) {
			return []string{obj.(string)}, nil
		},
		"type": func(obj interface{}) ([]string, error) {
			return []string{fmt.Sprintf("%T", obj)}, nil
		},
	}

	// Add multiple indexers
	err := si.addIndexers(newIndexers)
	assert.Nil(t, err)
}
