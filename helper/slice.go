package helper

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Filter filters values from a slice using a filter function.
// It returns a new slice with only the elements of s
// for which f returned true.
func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Merge - receives slices of type T and merges them into a single slice of type T.
func Merge[T any](slices ...[]T) (mergedSlice []T) {
	for _, slice := range slices {
		mergedSlice = append(mergedSlice, slice...)
	}
	return mergedSlice
}

// Includes - given a slice of type T and a value of type T,
// determines whether the value is contained by the slice.
func Includes[T comparable](slice []T, value T) bool {
	for _, el := range slice {
		if el == value {
			return true
		}
	}
	return false
}

// Sort - sorts given a slice of any order-able type T
// The constraints.Ordered constraint in the Sort() function guarantees that
// the function can sort values of any type supporting the operators <, <=, >=, >.
func Sort[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}
