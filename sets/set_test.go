package sets

import (
	"reflect"
	"testing"
)

func TestIntSet(t *testing.T) {
	// Test NewSet and Insert with int
	set := NewSet(1, 2, 3)
	if !set.Has(1) || !set.Has(2) || !set.Has(3) {
		t.Errorf("NewSet or Insert with int failed")
	}

	// Test Delete with int
	set.Delete(1)
	if set.Has(1) {
		t.Errorf("Delete with int failed")
	}

	// Test HasAll with int
	if !set.HasAll(2, 3) {
		t.Errorf("HasAll with int failed")
	}

	// Test HasAny with int
	if !set.HasAny(2, 4) {
		t.Errorf("HasAny with int failed")
	}

	// Test Difference with int
	otherSet := NewSet(2, 4)
	diffSet := set.Difference(otherSet)
	if diffSet.Has(2) || !diffSet.Has(3) {
		t.Errorf("Difference with int failed")
	}

	// Test Union with int
	unionSet := set.Union(otherSet)
	if !unionSet.HasAll(2, 3, 4) {
		t.Errorf("Union with int failed")
	}

	// Test Intersection with int
	intersectSet := set.Intersection(otherSet)
	if !intersectSet.Has(2) || intersectSet.Has(3) {
		t.Errorf("Intersection with int failed")
	}

	// Test IsSuperset with int
	if !set.IsSuperset(NewSet(2)) {
		t.Errorf("IsSuperset with int failed")
	}

	// Test Equal with int
	if !set.Equal(NewSet(2, 3)) {
		t.Errorf("Equal with int failed")
	}

	// Test List with int
	list := set.List(func(lhs, rhs int) bool { return lhs < rhs })
	if len(list) != 2 || list[0] != 2 || list[1] != 3 {
		t.Errorf("List with int failed")
	}

	// Test UnsortedList with int
	unsortedList := set.UnsortedList()
	if len(unsortedList) != 2 {
		t.Errorf("UnsortedList with int failed")
	}

	// Test PopAny with int
	_, ok := set.PopAny()
	if !ok || len(set) != 1 {
		t.Errorf("PopAny with int failed")
	}

	// Test Len with int
	if set.Len() != 1 {
		t.Errorf("Len with int failed")
	}
}

func TestStringSet(t *testing.T) {
	s := Set[string]{}
	s2 := Set[string]{}
	if len(s) != 0 {
		t.Errorf("Expected len=0: %d", len(s))
	}
	s.Insert("a", "b")
	if len(s) != 2 {
		t.Errorf("Expected len=2: %d", len(s))
	}
	s.Insert("c")
	if s.Has("d") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Has("a") {
		t.Errorf("Missing contents: %#v", s)
	}
	s.Delete("a")
	if s.Has("a") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s.Insert("a")
	if s.HasAll("a", "b", "d") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.HasAll("a", "b") {
		t.Errorf("Missing contents: %#v", s)
	}
	s2.Insert("a", "b", "d")
	if s.IsSuperset(s2) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	s2.Delete("d")
	if !s.IsSuperset(s2) {
		t.Errorf("Missing contents: %#v", s)
	}
}

func TestStringSetDeleteMultiples(t *testing.T) {
	s := Set[string]{}
	s.Insert("a", "b", "c")
	if len(s) != 3 {
		t.Errorf("Expected len=3: %d", len(s))
	}

	s.Delete("a", "c")
	if len(s) != 1 {
		t.Errorf("Expected len=1: %d", len(s))
	}
	if s.Has("a") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if s.Has("c") {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Has("b") {
		t.Errorf("Missing contents: %#v", s)
	}

}

func TestNewStringSet(t *testing.T) {
	s := NewSet("a", "b", "c")
	if len(s) != 3 {
		t.Errorf("Expected len=3: %d", len(s))
	}
	if !s.Has("a") || !s.Has("b") || !s.Has("c") {
		t.Errorf("Unexpected contents: %#v", s)
	}
}

func TestStringSetList(t *testing.T) {
	s := NewSet("z", "y", "x", "a")
	if !reflect.DeepEqual(s.List(func(lhs, rhs string) bool {
		return lhs < rhs
	}), []string{"a", "x", "y", "z"}) {
		t.Errorf("List gave unexpected result: %#v", s.List(func(lhs, rhs string) bool {
			return lhs < rhs
		}))
	}
}

func TestStringSetDifference(t *testing.T) {
	a := NewSet("1", "2", "3")
	b := NewSet("1", "2", "4", "5")
	c := a.Difference(b)
	d := b.Difference(a)
	if len(c) != 1 {
		t.Errorf("Expected len=1: %d", len(c))
	}
	if !c.Has("3") {
		t.Errorf("Unexpected contents: %#v", c.List(func(lhs, rhs string) bool {
			return lhs < rhs
		}))
	}
	if len(d) != 2 {
		t.Errorf("Expected len=2: %d", len(d))
	}
	if !d.Has("4") || !d.Has("5") {
		t.Errorf("Unexpected contents: %#v", d.List(func(lhs, rhs string) bool {
			return lhs < rhs
		}))
	}
}

func TestStringSetHasAny(t *testing.T) {
	a := NewSet("1", "2", "3")

	if !a.HasAny("1", "4") {
		t.Errorf("expected true, got false")
	}

	if a.HasAny("0", "4") {
		t.Errorf("expected false, got true")
	}
}

func TestStringSetEquals(t *testing.T) {
	// Simple case (order doesn't matter)
	a := NewSet("1", "2")
	b := NewSet("2", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// It is a set; duplicates are ignored
	b = NewSet("2", "2", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	// Edge cases around empty sets / empty strings
	a = NewSet[string]()
	b = NewSet[string]()
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	b = NewSet("1", "2", "3")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	b = NewSet("1", "2", "")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	// Check for equality after mutation
	a = NewSet[string]()
	a.Insert("1")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert("2")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}

	a.Insert("")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	a.Delete("")
	if a.Equal(b) {
		t.Errorf("Expected to be not-equal: %v vs %v", a, b)
	}
}

func TestStringUnion(t *testing.T) {
	tests := []struct {
		s1       Set[string]
		s2       Set[string]
		expected Set[string]
	}{
		{
			NewSet("1", "2", "3", "4"),
			NewSet("3", "4", "5", "6"),
			NewSet("1", "2", "3", "4", "5", "6"),
		},
		{
			NewSet("1", "2", "3", "4"),
			NewSet[string](),
			NewSet("1", "2", "3", "4"),
		},
		{
			NewSet[string](),
			NewSet("1", "2", "3", "4"),
			NewSet("1", "2", "3", "4"),
		},
		{
			NewSet[string](),
			NewSet[string](),
			NewSet[string](),
		},
	}

	for _, test := range tests {
		union := test.s1.Union(test.s2)
		if union.Len() != test.expected.Len() {
			t.Errorf("Expected union.Len()=%d but got %d", test.expected.Len(), union.Len())
		}

		if !union.Equal(test.expected) {
			t.Errorf("Expected union.Equal(expected) but not true.  union:%v expected:%v", union.List(func(lhs, rhs string) bool {
				return lhs < rhs
			}), test.expected.List(func(lhs, rhs string) bool {
				return lhs < rhs
			}))
		}
	}
}

func TestStringIntersection(t *testing.T) {
	tests := []struct {
		s1       Set[string]
		s2       Set[string]
		expected Set[string]
	}{
		{
			NewSet("1", "2", "3", "4"),
			NewSet("3", "4", "5", "6"),
			NewSet("3", "4"),
		},
		{
			NewSet("1", "2", "3", "4"),
			NewSet("1", "2", "3", "4"),
			NewSet("1", "2", "3", "4"),
		},
		{
			NewSet("1", "2", "3", "4"),
			NewSet[string](),
			NewSet[string](),
		},
		{
			NewSet[string](),
			NewSet("1", "2", "3", "4"),
			NewSet[string](),
		},
		{
			NewSet[string](),
			NewSet[string](),
			NewSet[string](),
		},
	}

	for _, test := range tests {
		intersection := test.s1.Intersection(test.s2)
		if intersection.Len() != test.expected.Len() {
			t.Errorf("Expected intersection.Len()=%d but got %d", test.expected.Len(), intersection.Len())
		}

		if !intersection.Equal(test.expected) {
			t.Errorf("Expected intersection.Equal(expected) but not true.  intersection:%v expected:%v", intersection.List(func(lhs, rhs string) bool {
				return lhs < rhs
			}), test.expected.List(func(lhs, rhs string) bool {
				return lhs < rhs
			}))
		}
	}
}
