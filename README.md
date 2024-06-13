# Cache Package

The `cache` package provides a thread-safe, indexed cache implementation, and flexible caching mechanism with various eviction policies such as FIFO (First In, First Out), LRU (Least Recently Used), and LFU (Least Frequently Used). It offers the following features:

- **Multiple Eviction Policies**: Supports FIFO, LRU, and LFU eviction policies.

- **Thread-Safe Operations**: The cache operations are designed to be safe for concurrent use, allowing multiple goroutines to read and write to the cache concurrently.

- **Indexing**: The cache supports indexing objects based on their properties. This allows for efficient lookup of objects based on specific criteria.

- **Custom Key Generation**: Users can provide a custom key generation function to generate keys for objects stored in the cache.

# Installation

To use the `cache` package in your Go project, you can use the `go get` command:

```bash
go get -u github.com/liuxinbot/cache
```

# Usage
## Creating a Cache with Custom Key Generation
You can create a cache with custom key generation using the NewStore function:
```go
package main

import (
	"fmt"

	"github.com/liuxinbot/cache"
)

func main() {
	// Define a custom key generation function
	keyFunc := func(obj interface{}) (string, error) {
		return obj.(string), nil
	}

	// Create a new cache with the custom key generation function
	store := cache.NewStore(keyFunc)

	// Add objects to the cache
	store.Add("apple")
	store.Add("banana")
	store.Add("orange")

	// Retrieve an object by key
	item, exists, err := store.Get("banana")
	if err != nil {
		fmt.Println("Error getting object:", err)
		return
	}

	if exists {
		fmt.Println("Found object:", item)
	} else {
		fmt.Println("Object not found")
	}
}
```

## Creating a Cache with Indexing
You can create a cache with indexing using the NewIndexer function:

```go
package main

import (
	"fmt"

	"github.com/liuxinbot/cache"
)

func main() {
	// Define a custom key generation function
	keyFunc := func(obj interface{}) (string, error) {
		return obj.(string), nil
	}

	// Define an indexer function
	indexer := func(obj interface{}) ([]any, error) {
		return []any{len(obj.(string))}, nil
	}

	// Create a new indexed cache
	indexedStore := cache.NewIndexer[any](keyFunc)

	// add indexer
	indexers := cache.Indexers[any]{"length": indexer}
	indexedStore.AddIndexers(indexers)

	// Add objects to the cache
	indexedStore.Add("apple")
	indexedStore.Add("banana")
	indexedStore.Add("orange")

	// Retrieve objects by index
	items, err := indexedStore.ListByIndex("length", 5)
	if err != nil {
		fmt.Println("Error listing objects by index:", err)
		return
	}

	// Print the retrieved objects
	for _, item := range items {
		fmt.Println(item)
	}
}

```

## Creating an Eviction Cache
You can create a new eviction cache by specifying the key function, eviction policy, and indexers.

```go
package main

import (
	"fmt"

	"github.com/liuxinbot/cache"
	"github.com/liuxinbot/cache/eviction"
)

// Key function for the cache
func keyFunc(obj interface{}) (int, error) {
	return obj.(int), nil
}

func main() {
	// Create a new FIFO eviction cache with a capacity of 2
	fifoPolicy := eviction.NewFIFO[int](2)
	store := cache.NewEvictionCache(keyFunc, fifoPolicy, make(cache.Indexers[int]))

	// Add items to the cache
	store.Add(1)
	store.Add(2)

	// Print the current size of the cache
	fmt.Println("Cache Size:", store.Size())

	// Add another item, causing eviction
	store.Add(3)

	// Print the current size of the cache after eviction
	fmt.Println("Cache Size after eviction:", store.Size())
}
```

### Eviction Policies
#### FIFO (First In, First Out)
```go
fifoPolicy := eviction.NewFIFO[int](capacity)
cache := cache.NewEvictionCache(keyFunc, fifoPolicy, make(cache.Indexers[int]))
```

#### LRU (Least Recently Used)
```go
lruPolicy := eviction.NewLRU[int](capacity)
cache := cache.NewEvictionCache(keyFunc, lruPolicy, make(cache.Indexers[int]))
```

#### LFU (Least Frequently Used)
```go
lfuPolicy := eviction.NewLFU[int](capacity)
cache := cache.NewEvictionCache(keyFunc, lfuPolicy, make(cache.Indexers[int]))
```


# Testing
The cache package includes comprehensive unit tests to ensure the correctness of its functionality. You can run the tests using the go test command:

```bash
go test -v github.com/liuxinbot/cache
```

# License
This project is licensed under the Apache-2.0 License. See the LICENSE file for details.

# Contributing
Contributions are welcome! Please open an issue or submit a pull request with your changes. For major changes, please open an issue first to discuss what you would like to change.

# Contact
For any questions or suggestions, feel free to contact the project maintainers.
