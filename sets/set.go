package sets

import (
	"reflect"
	"sort"
)

type Empty struct{}

// Set is a generic set implemented via map for minimal memory consumption.
type Set[T comparable] map[T]Empty

// NewSet creates a Set from a list of values.
func NewSet[T comparable](items ...T) Set[T] {
	ss := Set[T]{}
	ss.Insert(items...)
	return ss
}

// KeySet creates a Set from the keys of a map[K](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func KeySet[K comparable](theMap interface{}) Set[K] {
	v := reflect.ValueOf(theMap)
	ret := Set[K]{}

	for _, keyValue := range v.MapKeys() {
		ret.Insert(keyValue.Interface().(K))
	}
	return ret
}

// Insert adds items to the set.
func (s Set[T]) Insert(items ...T) Set[T] {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

// Delete removes all items from the set.
func (s Set[T]) Delete(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if and only if item is contained in the set.
func (s Set[T]) Has(item T) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true if and only if all items are contained in the set.
func (s Set[T]) HasAll(items ...T) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (s Set[T]) HasAny(items ...T) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

// Difference returns a set of objects that are not in s2.
func (s Set[T]) Difference(s2 Set[T]) Set[T] {
	result := NewSet[T]()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
func (s Set[T]) Union(s2 Set[T]) Set[T] {
	result := NewSet[T]()
	for key := range s {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2.
func (s Set[T]) Intersection(s2 Set[T]) Set[T] {
	var walk, other Set[T]
	result := NewSet[T]()
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Set[T]) IsSuperset(s2 Set[T]) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
func (s Set[T]) Equal(s2 Set[T]) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

// List returns the contents as a sorted slice.
func (s Set[T]) List(less func(lhs, rhs T) bool) []T {
	res := make([]T, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Slice(res, func(i, j int) bool {
		return less(res[i], res[j])
	})
	return res
}

// UnsortedList returns the slice with contents in random order.
func (s Set[T]) UnsortedList() []T {
	res := make([]T, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// PopAny returns a single element from the set.
func (s Set[T]) PopAny() (T, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}
	var zeroValue T
	return zeroValue, false
}

// Len returns the size of the set.
func (s Set[T]) Len() int {
	return len(s)
}
