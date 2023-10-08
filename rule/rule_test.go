package rule_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sqjian/go-kit/rule"
	"reflect"
	"testing"
	"time"
)

func TestExecRule(t *testing.T) {
	type args struct {
		code string
		env  any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				code: `(get(fromJSON(ava_slot),"age")=="青年")?"年轻人":"老头子"`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    "年轻人",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				code: `(get(fromJSON(ava_slot),"age")=="青年")/*解析原始json*/?/*三目运算符*/"年轻人":"老头子"`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    "年轻人",
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				code: `(get(fromJSON(ava_slot),"age")=="青年")?"年轻人":"老头子"//测试用`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    "年轻人",
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				code: `(get(fromJSON(ava_slot),"age")=="青年")?"年轻人":"老头子"//测试用`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    "年轻人",
			wantErr: false,
		},
		{
			name: "test5",
			args: args{
				code: `true`,
				env:  nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test6",
			args: args{
				code: `6`,
				env:  nil,
			},
			want:    6,
			wantErr: false,
		},
		{
			name: "test7",
			args: args{
				code: `[1,2,3]`,
				env:  nil,
			},
			want:    []any{1, 2, 3},
			wantErr: false,
		},
		{
			name: "test8",
			args: args{
				code: `{a: 1, b: 2, c: 3}`,
				env:  nil,
			},
			want:    map[string]any{"a": 1, "b": 2, "c": 3},
			wantErr: false,
		},
		{
			name: "test9",
			args: args{
				code: `nil`,
				env:  nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "test10",
			args: args{
				code: `(1+1)/2`,
				env:  nil,
			},
			want:    float64(1),
			wantErr: false,
		},
		{
			name: "test11",
			args: args{
				code: `true?1:2`,
				env:  nil,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "test12",
			args: args{
				code: `false?1:2`,
				env:  nil,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "test13",
			args: args{
				code: `a?.B??"xxx"`,
				env: map[string]any{
					"a": struct {
						B any
					}{B: nil}},
			},
			want:    "xxx",
			wantErr: false,
		},
		{
			name: "test14",
			args: args{
				code: `fromJSON(ava_slot)["age"]`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    "青年",
			wantErr: false,
		},
		{
			name: "test15",
			args: args{
				code: `"age" in fromJSON(ava_slot)`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test16",
			args: args{
				code: `get(fromJSON(ava_slot),"age")`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    "青年",
			wantErr: false,
		},
		{
			name: "test17",
			args: args{
				code: `get(fromJSON(ava_slot),"age")+"xxx"`,
				env: map[string]any{
					"ava_slot": `{"age":"青年","appearance":"五官端正","character":"阳光温暖","field":"融媒体","gender":"女"}`,
				},
			},
			want:    "青年xxx",
			wantErr: false,
		},
		{
			name: "test18",
			args: args{
				code: `"666" matches "\\d{3}"`,
				env:  nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test19",
			args: args{
				code: `Age in 18..45 and Name not in ["admin", "root"]`,
				env: map[string]any{
					"Age":  21,
					"Name": "21",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test20",
			args: args{
				code: `"a"|upper()`,
				env:  nil,
			},
			want:    "A",
			wantErr: false,
		},
		{
			name: "test21",
			args: args{
				code: `date("2023-08-14") + duration("1h")`,
				env:  nil,
			},
			want: func() time.Time {
				t, _ := time.Parse("2006-01-02T15:04:05Z", "2023-08-14T01:00:00Z")
				return t
			}(),
			wantErr: false,
		},
		{
			name: "test22",
			args: args{
				code: `all([1,2,3],#>0)`,
				env:  nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test23",
			args: args{
				code: `any([1,2,3],#>0)`,
				env:  nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test24",
			args: args{
				code: `one([1,2,3],#>0)`,
				env:  nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "test25",
			args: args{
				code: `none([1,2,3],#>0)`,
				env:  nil,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "test26",
			args: args{
				code: `map([1,2,3],#+1)`,
				env:  nil,
			},
			want:    []any{2, 3, 4},
			wantErr: false,
		},
		{
			name: "test27",
			args: args{
				code: `filter([1,2,3],#>2)`,
				env:  nil,
			},
			want:    []any{3},
			wantErr: false,
		},
		{
			name: "test28",
			args: args{
				code: `findIndex([1, 2, 3, 4], # > 20)`,
				env:  nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "test29",
			args: args{
				code: `groupBy([1,2,3], #>2)`,
				env:  nil,
			},
			want: map[any][]any{
				any(false): {1, 2},
				any(true):  {3},
			},
			wantErr: false,
		},
		{
			name: "test30",
			args: args{
				code: `reduce(1..9, #acc + #, 1)`,
				env:  nil,
			},
			want:    46,
			wantErr: false,
		},
		{
			name: "test31",
			args: args{
				code: `type("hello")`,
				env:  nil,
			},
			want:    "string",
			wantErr: false,
		},
		{
			name: "test32",
			args: args{
				code: `trim("__Hello__", "_")`,
				env:  nil,
			},
			want:    "Hello",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.ExecRule(tt.args.code, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecRule() got = %v, want %v", spew.Sdump(got), spew.Sdump(tt.want))
			}
		})
	}
}

func TestExecRuleRet(t *testing.T) {
	type args struct {
		code string
		env  any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{

		{
			name: "test32",
			args: args{
				code: `1+1`,
				env:  nil,
			},
			want:    "Hello",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.ExecRule(tt.args.code, tt.args.env, rule.AsString())
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecRule() got = %v, want %v", spew.Sdump(got), spew.Sdump(tt.want))
			}
		})
	}
}
