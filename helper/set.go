package helper

// Set is a set of any values.
type Set[T comparable] map[T]struct{}

// MakeSet returns a set of some element type.
func MakeSet[T comparable]() Set[T] {
	return make(Set[T])
}

// Add adds v to the set s.
// If v is already in s this has no effect.
func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

// Delete removes v from the set s.
// If v is not in s this has no effect.
func (s Set[T]) Delete(v T) {
	delete(s, v)
}

// Contains reports whether v is in s.
func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// Len reports the number of elements in s.
func (s Set[T]) Len() int {
	return len(s)
}

// All return all elements.
func (s Set[T]) All() []T {
	var _t []T
	for i := range s {
		_t = append(_t, i)
	}
	return _t
}
