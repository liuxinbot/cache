package cache

import "fmt"

// Store defines a basic storage interface.
type Store[T comparable] interface {
	// Add inserts an object.
	Add(obj interface{}) error

	// Update modifies an existing object.
	Update(obj interface{}) error

	// Delete removes an object.
	Delete(obj interface{}) error

	// List returns all objects.
	List() []interface{}

	// ListKeys returns all keys.
	ListKeys() []T

	// Get returns an object by its key.
	Get(obj interface{}) (interface{}, bool, error)

	// GetByKey returns an object by its key string.
	GetByKey(key T) (interface{}, bool, error)

	// Replace replaces all objects with the given list.
	Replace([]interface{}) error

	// Size returns count of object.
	Size() int
}

// KeyFunc generates a key from an object.
type KeyFunc[T comparable] func(obj interface{}) (T, error)

// KeyError represents an error during key generation.
type KeyError struct {
	Obj interface{}
	Err error
}

// Error returns a human-readable description of the KeyError.
func (k KeyError) Error() string {
	return fmt.Sprintf("couldn't create key for object %+v: %v", k.Obj, k.Err)
}

// Unwrap returns the underlying error.
func (k KeyError) Unwrap() error {
	return k.Err
}
