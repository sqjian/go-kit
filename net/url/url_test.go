package url_test

import (
	"github.com/sqjian/go-kit/net/url"
	"testing"
)

func TestIsValidUrl(t *testing.T) {
	type args struct {
		u1 string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{u1: "htt://x.x.x"}, wantErr: true},
		{name: "test2", args: args{u1: "http://x.x.x"}, wantErr: false},
		{name: "test3", args: args{u1: "x.x.x"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := url.IsValidUrl(tt.args.u1); (err != nil) != tt.wantErr {
				t.Errorf("IsValidUrl() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("IsValidUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckUrl(t *testing.T) {
	type args[T interface{ []byte | string }] struct {
		u T
	}
	type testCase[T interface{ []byte | string }] struct {
		name    string
		args    args[T]
		wantErr bool
	}
	tests := []testCase[[]byte]{
		{
			name: "test1",
			args: args[[]byte]{
				u: []byte("http://x.x.x"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := url.CheckUrl(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("CheckUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
