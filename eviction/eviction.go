package eviction

// Policy defines the interface for cache eviction policies.
type Policy[T comparable] interface {
	Put(key T) (T, bool) // Adds a key to the cache, returns the evicted key if any.
	Delete(key T)        // Removes a key from the cache.
	Evict() (T, bool)    // Evicts a key from the cache based on the policy.
	Reset()              // Clears all keys from the cache.
	Size() int           // Returns the current number of keys in the cache.
}
