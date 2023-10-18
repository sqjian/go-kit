package json_test

import (
	_ "embed"
	"encoding/json"
	easyjson "github.com/sqjian/go-kit/encoding/json"
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
			gotValue, err := easyjson.Set(tt.args.data, tt.args.setValue, tt.args.keys...)
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
		}, {
			name: "test2",
			args: args{
				data: []byte(`{"age": "儿童","appearance": "五官端正","character": ">甜美可爱","field": "自媒体","gender": "女"}`),
				keys: []string{"age"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotDataType, gotOffset, err := easyjson.Get(tt.args.data, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("gotValue:%v,gotDataType:%v,gotOffset:%v", string(gotValue), gotDataType.String(), gotOffset)
		})
	}
}

//go:embed testdata/person.jsonc
var personCommented []byte

func TestStandardize(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "person.jsonc",
			args: args{
				data: personCommented,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := easyjson.Standardize(tt.args.data)
			v := make(map[string]any)
			err := json.Unmarshal(got, &v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Standardize() decode = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(v)
		})
	}
}
