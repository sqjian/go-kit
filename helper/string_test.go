package helper

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"
)

func TestSplitAfter(t *testing.T) {
	type args struct {
		str string
		sep []rune
		min int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				str: `中秋佳节，月圆人团圆，愿老师身体健康，工作顺利。`,
				sep: []rune{'，', '。'},
				min: 10,
			},
			want:    []string{"中秋佳节，月圆人团圆，", "愿老师身体健康，工作顺利。"},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				str: `中秋佳节，月圆人团圆，。愿老师身体健康，！工作顺利。`,
				sep: []rune{'，', '。', '!'},
				min: 10,
			},
			want:    []string{"中秋佳节，月圆人团圆，", "。愿老师身体健康，！工作顺利。"},
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				str: `中秋佳节，月圆人团圆，。愿老师身体健康，！工作顺利。。`,
				sep: []rune{'，', '。', '!'},
				min: 10,
			},
			want:    []string{"中秋佳节，月圆人团圆，", "。愿老师身体健康，！工作顺利。。"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitAfter(tt.args.str, tt.args.sep, tt.args.min)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitAfter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitAfter() got = %v, want %v", got, tt.want)
			}
			spew.Dump(got)
		})
	}
}
