package helper

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IsValidUrl(tt.args.u1); (err != nil) != tt.wantErr {
				t.Errorf("IsValidUrl() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("IsValidUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
