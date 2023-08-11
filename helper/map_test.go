package helper

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"testing"
)

func TestKeys(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []K
	}
	tests := []testCase[string, int]{
		{
			name: "test1",
			args: args[string, int]{
				m: map[string]int{"name": 1},
			},
			want: []string{"name"},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Keys(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	type args[K comparable, V interface {
		constraints.Float | constraints.Integer
	}] struct {
		m map[K]V
	}
	type testCase[K comparable, V interface {
		constraints.Float | constraints.Integer
	}] struct {
		name string
		args args[K, V]
		want V
	}
	tests := []testCase[string, int]{
		{
			name: "test1",
			args: args[string, int]{
				m: map[string]int{"name": 1},
			},
			want: 1,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.m); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
