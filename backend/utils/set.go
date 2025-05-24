// utils/set.go
package utils

// Generic Set type
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	s := make(Set[T])
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Has(item T) bool {
	_, exists := s[item]
	return exists
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) ToSlice() []T {
	items := make([]T, 0, len(s))
	for item := range s {
		items = append(items, item)
	}
	return items
}
