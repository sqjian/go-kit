package json_test

import (
	"bytes"
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
			if !bytes.Contains(gotValue, tt.args.setValue) {
				t.Fatalf("set failed.")
			}
			t.Logf("raw:%v,processed:%v", string(tt.args.data), string(gotValue))
		})
	}
}

func TestMustGet(t *testing.T) {
	type args struct {
		data []byte
		keys []string
	}
	tests := []struct {
		name   string
		args   args
		expect []byte
	}{
		{
			name: "test1",
			args: args{
				data: []byte(`{"name":"John","age":30,"city":"New York"}`),
				keys: []string{"city"},
			},
			expect: []byte("New York"),
		}, {
			name: "test2",
			args: args{
				data: []byte(`{"age": "儿童","appearance": "五官端正","character": ">甜美可爱","field": "自媒体","gender": "女"}`),
				keys: []string{"age"},
			},
			expect: []byte("儿童"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue := easyjson.MustGet(tt.args.data, tt.args.keys...)
			if string(gotValue) != string(tt.expect) {
				t.Errorf("MustGet() failed, gotVal:%s,expect:%s", gotValue, tt.expect)
				return
			}
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
		wantVal []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				data: []byte(`{"name":"John","age":30,"city":"New York"}`),
				keys: []string{"city"},
			},
			wantVal: []byte("New York"),
			wantErr: false,
		}, {
			name: "test2",
			args: args{
				data: []byte(`{"age": "儿童","appearance": "五官端正","character": ">甜美可爱","field": "自媒体","gender": "女"}`),
				keys: []string{"age"},
			},
			wantVal: []byte("儿童"),
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
			if string(gotValue) != string(tt.wantVal) {
				t.Errorf("MustGet() failed, gotVal:%s,wantVal:%s", gotValue, tt.wantVal)
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
			if bytes.Contains(got, []byte("//")) {
				t.Fatalf("still containing comment chars")
			}
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

func TestMarshal(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				data: map[string]string{"key": "<val>"},
			},
			want:    []byte(`{"key":"<val>"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := easyjson.Marshal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bytes.Contains(got, []byte("<")) || bytes.Contains(got, []byte(">")) {
				return
			} else {
				t.Logf("do not contains < or >")
			}
		})
	}
}
