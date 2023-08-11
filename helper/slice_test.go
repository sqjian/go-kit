package helper_test

import (
	"github.com/sqjian/go-kit/helper"
	"golang.org/x/exp/constraints"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	type args[T any] struct {
		s []T
		f func(T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{
				s: []int{1, 2, 3},
				f: func(i int) bool {
					if i < 2 {
						return false
					}
					return true
				},
			},
			want: []int{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.Filter(tt.args.s, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	type args[T comparable] struct {
		slice []T
		value T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{
				slice: []int{1, 2, 3},
				value: 2,
			},
			want: true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.Includes(tt.args.slice, tt.args.value); got != tt.want {
				t.Errorf("Includes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	type args[T any] struct {
		slices [][]T
	}
	type testCase[T any] struct {
		name            string
		args            args[T]
		wantMergedSlice []T
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{
				slices: [][]int{{1, 2, 3}, {4, 5, 6}},
			},
			wantMergedSlice: []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMergedSlice := helper.Merge(tt.args.slices...); !reflect.DeepEqual(gotMergedSlice, tt.wantMergedSlice) {
				t.Errorf("Merge() = %v, want %v", gotMergedSlice, tt.wantMergedSlice)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type args[T constraints.Ordered] struct {
		s []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{
				s: []int{1, 2, 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper.Sort(tt.args.s)
		})
	}
}
