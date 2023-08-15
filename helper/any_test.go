package helper_test

import (
	"github.com/sqjian/go-kit/helper"
	"testing"
)

func TestIsPtr(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test1", args: args{v: new(int)}, want: true},
		{name: "test2", args: args{v: 1}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.IsPtr(tt.args.v); got != tt.want {
				t.Errorf("IsPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}
