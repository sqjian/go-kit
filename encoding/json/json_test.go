package json_test

import (
	"github.com/sqjian/go-kit/encoding/json"
	"testing"
)

func TestSet(t *testing.T) {
	type args struct {
		data     []byte
		setValue []byte
		keys     []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				data:     []byte(`{"name":"John","age":30,"city":"New York"}`),
				setValue: []byte(`"San Francisco"`),
				keys:     []string{"city"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := json.Set(tt.args.data, tt.args.setValue, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("raw:%v,processed:%v", string(tt.args.data), string(gotValue))
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		data []byte
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				data: []byte(`{"name":"John","age":30,"city":"New York"}`),
				keys: []string{"city"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotDataType, gotOffset, err := json.Get(tt.args.data, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("gotValue:%v,gotDataType:%v,gotOffset:%v", string(gotValue), gotDataType.String(), gotOffset)
		})
	}
}
