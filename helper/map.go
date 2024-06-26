package helper

import (
	"golang.org/x/exp/constraints"
	"sort"
)

// Keys returns the keys of the map m in a slice.
// The keys will be returned in an unpredictable order.
// This function has two type parameters, K and V.
// Map keys must be comparable, so key has the predeclared
// constraint comparable. Map values can be any type.
func Keys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// Sum sums the values of map containing numeric or float values
func Sum[K comparable, V constraints.Float | constraints.Integer](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func SortedMap[K comparable, V any](m map[K]V, less func(a, b K) bool, f func(k K, v V)) {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return less(keys[i], keys[j])
	})
	for _, k := range keys {
		f(k, m[k])
	}
}
