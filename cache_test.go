package cache

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Key function for testing
func testKeyFunc(obj interface{}) (string, error) {
	return obj.(string), nil
}

func TestCache(t *testing.T) {
	store := NewStore(testKeyFunc)

	// Test Add
	err := store.Add("test1")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(store.List()))

	// Test Get
	item, exists, err := store.Get("test1")
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, "test1", item)

	// Test Update
	err = store.Update("test2")
	assert.Nil(t, err)
	item, exists, err = store.Get("test2")
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, "test2", item)

	// Test Delete
	err = store.Delete("test2")
	assert.Nil(t, err)
	_, exists, err = store.Get("test2")
	assert.Nil(t, err)
	assert.False(t, exists)

	// Test Replace
	newItems := []interface{}{"new1", "new2", "new3"}
	err = store.Replace(newItems)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(store.List()))
}

func TestIndexer(t *testing.T) {
	indexers := Indexers[any]{
		"name": func(obj interface{}) ([]any, error) {
			return []any{obj.(string)}, nil
		},
	}
	store := NewIndexer[any](testKeyFunc)
	store.AddIndexers(indexers)

	// Test AddIndexer
	err := store.AddIndexer("name2", func(obj interface{}) ([]any, error) {
		return []any{obj.(string)}, nil
	})
	assert.Nil(t, err)

	// Test Add
	err = store.Add("test1")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(store.List()))

	// Test ListByIndex
	items, err := store.ListByIndex("name", "test1")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "test1", items[0])
}

// Benchmark testing
func BenchmarkCacheAdd(b *testing.B) {
	store := NewStore(testKeyFunc)

	for i := 0; i < b.N; i++ {
		store.Add("test" + fmt.Sprintf("%d", i))
	}
}

func BenchmarkCacheGet(b *testing.B) {
	store := NewStore(testKeyFunc)
	store.Add("test1")

	for i := 0; i < b.N; i++ {
		store.Get("test1")
	}
}

func BenchmarkCacheUpdate(b *testing.B) {
	store := NewStore(testKeyFunc)
	store.Add("test1")

	for i := 0; i < b.N; i++ {
		store.Update("test1")
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	store := NewStore(testKeyFunc)
	store.Add("test1")

	for i := 0; i < b.N; i++ {
		store.Delete("test1")
	}
}

// Example test
func ExampleCache() {
	store := NewStore(testKeyFunc)

	// Add items to the cache
	store.Add("item1")
	store.Add("item2")

	// List all items
	items := store.List()
	sort.Slice(items, func(i, j int) bool {
		return items[i].(string) < items[j].(string)
	})
	fmt.Println(items)

	// Get an item by key
	item, exists, _ := store.Get("item1")
	if exists {
		fmt.Println("Got:", item)
	}

	// Update an item
	store.Update("item1_updated")
	item, exists, _ = store.Get("item1_updated")
	if exists {
		fmt.Println("Updated to:", item)
	}

	// Delete an item
	store.Delete("item2")
	_, exists, _ = store.Get("item2")
	if !exists {
		fmt.Println("Deleted item2")
	}

	// Output:
	// [item1 item2]
	// Got: item1
	// Updated to: item1_updated
	// Deleted item2
}
