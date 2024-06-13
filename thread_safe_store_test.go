package cache

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMultiTypeIndexValue tests the ThreadSafeStore with multiple index value types.
func TestMultiTypeIndexValue(t *testing.T) {
	// User represents a sample struct for testing.
	type User struct {
		Name string
		Age  int
		Sex  string
	}

	// Initialize the store with empty indexers and indices.
	store := NewThreadSafeStore[any, int](Indexers[any]{}, Indexes[any, int]{})

	// Add "age" indexer.
	err := store.AddIndexer("age", func(obj any) ([]any, error) {
		u, ok := obj.(*User)
		if !ok {
			return nil, fmt.Errorf("object is not a *User")
		}
		return []any{u.Age}, nil
	})
	if err != nil {
		t.Fatalf("failed to add age index: %v", err)
	}

	// Add "sex" indexer.
	err = store.AddIndexers(Indexers[any]{
		"sex": func(obj any) ([]any, error) {
			u, ok := obj.(*User)
			if !ok {
				return nil, fmt.Errorf("object is not a *User")
			}
			return []any{u.Sex}, nil
		},
	})
	if err != nil {
		t.Fatalf("failed to add sex index: %v", err)
	}

	// Add test data to the store.
	for i := 0; i < 10; i++ {
		age := 10
		sex := "man"
		if i%2 == 0 {
			sex = "woman"
			age = 20
		}
		store.Add(i, &User{
			Name: fmt.Sprintf("name-%d", i),
			Age:  age,
			Sex:  sex,
		})
	}

	// Test querying by "age" index.
	res, err := store.IndexKeys("age", 10, func(lhs, rhs int) bool {
		return lhs < rhs
	})
	if err != nil {
		t.Fatalf("failed to query by age index: %v", err)
	}

	if !reflect.DeepEqual(res, []int{1, 3, 5, 7, 9}) {
		t.Errorf("expected [1,3,5,7,9], got %v", res)
	}

	// Test querying by "sex" index.
	res, err = store.IndexKeys("sex", "woman", func(lhs, rhs int) bool {
		return lhs < rhs
	})
	if err != nil {
		t.Fatalf("failed to query by sex index: %v", err)
	}

	if !reflect.DeepEqual(res, []int{0, 2, 4, 6, 8}) {
		t.Errorf("expected [0,2,4,6,8], got %v", res)
	}

}

func TestThreadSafeStore(t *testing.T) {
	indexByLength := func(obj any) ([]string, error) {
		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("object is not a string")
		}
		return []string{strconv.Itoa(len(str))}, nil
	}

	indexByPrefix := func(obj any) ([]string, error) {
		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("object is not a string")
		}
		if len(str) < 2 {
			return []string{str}, nil
		}
		return []string{str[:2]}, nil
	}

	indexers := Indexers[string]{
		"length": indexByLength,
		"prefix": indexByPrefix,
	}
	indices := Indexes[string, string]{}

	store := NewThreadSafeStore[string, string](indexers, indices)

	store.Add("key1", "hello")
	store.Add("key2", "world")
	store.Add("key3", "he")
	store.Add("key4", "wo")
	store.Add("key5", "g")

	// Test Get
	item, exists := store.Get("key1")
	if !exists || item != "hello" {
		t.Errorf("expected 'hello', got %v", item)
	}

	// Test List
	expectedList := []any{"hello", "world", "he", "wo", "g"}
	assert.ElementsMatch(t, store.List(), expectedList)

	// Test ListKeys
	expectedKeys := []string{"key1", "key2", "key3", "key4", "key5"}
	assert.ElementsMatch(t, store.ListKeys(), expectedKeys)

	// Test Index
	indexedItems, err := store.Index("length", "hello", nil)
	assert.ElementsMatch(t, indexedItems, []any{"hello", "world"})

	indexedItems, err = store.Index("prefix", "he", nil)
	assert.ElementsMatch(t, indexedItems, []any{"hello", "he"})

	// Test ByIndex
	indexedItems, err = store.ByIndex("length", "5", nil)
	assert.ElementsMatch(t, indexedItems, []any{"hello", "world"})

	indexedItems, err = store.ByIndex("prefix", "he", nil)
	assert.ElementsMatch(t, indexedItems, []any{"hello", "he"})

	// Test Delete
	store.Delete("key1")
	item, exists = store.Get("key1")
	if exists {
		t.Errorf("expected 'key1' to be deleted")
	}

	expectedList = []any{"world", "he", "wo", "g"}
	assert.ElementsMatch(t, store.List(), expectedList)

	// Test Replace
	newItems := map[string]any{
		"key6": "new1",
		"key7": "new2",
	}
	store.Replace(newItems)

	expectedList = []any{"new1", "new2"}
	assert.ElementsMatch(t, store.List(), expectedList)
	expectedKeys = []string{"key6", "key7"}
	assert.ElementsMatch(t, store.ListKeys(), expectedKeys)

	// Test AddIndexers
	newIndexer := Indexers[string]{
		"suffix": func(obj any) ([]string, error) {
			str, ok := obj.(string)
			if !ok {
				return nil, fmt.Errorf("object is not a string")
			}
			if len(str) < 2 {
				return []string{str}, nil
			}
			return []string{str[len(str)-2:]}, nil
		},
	}
	err = store.AddIndexers(newIndexer)
	if err != nil {
		t.Errorf("unexpected error adding indexers: %v", err)
	}

	// Test ByIndex with new indexer
	store.Add("key8", "suffixTest")
	indexedItems, err = store.ByIndex("suffix", "st", nil)
	assert.ElementsMatch(t, indexedItems, []any{"suffixTest"})
}
