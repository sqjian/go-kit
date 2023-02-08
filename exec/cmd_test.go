package exec_test

import (
	"github.com/sqjian/go-kit/exec"
	"io"
	"os"
	"testing"
)

func TestCmd(t *testing.T) {
	type args struct {
		execName string
		opts     []exec.OptionFunc
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{execName: "ping", opts: []exec.OptionFunc{exec.WithArgs("www.baidu.com")}},
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{execName: "ping", opts: []exec.OptionFunc{exec.WithArgs("www.baidu.com"), exec.WithWriters(func() io.WriteCloser { f, _ := os.Create("testing.log"); return f }())}},
			wantErr: false,
		},
		{
			name:    "test3",
			args:    args{execName: "ping", opts: []exec.OptionFunc{exec.WithArgs("www.baidu.com"), exec.WithWriters(os.Stderr)}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := exec.Cmd(tt.args.execName, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Cmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
