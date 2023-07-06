package pool

import (
	"context"
	"testing"
)

func Test_genUniqueId(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test1", "->1"},
		{"test2", "->2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genUniqueId(context.Background()); got != tt.want {
				t.Errorf("genUniqueId() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
