package ctr

import (
	"testing"
)

func TestAes(t *testing.T) {
	type args struct {
		plainText []byte
		key       []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "test1", args: args{plainText: []byte("hello world"), key: []byte("0123456789012345")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("plainText：%s,key:%s", tt.args.plainText, tt.args.key)
			cipherText, err := AesEncrypt(tt.args.plainText, tt.args.key)
			if err != nil {
				t.Fatalf("AesEncrypt failed,err:%v", err)
			}
			t.Logf("cipherText data:%X", cipherText)
			decryptCode, err := AesDecrypt(cipherText, tt.args.key)
			if err != nil {
				t.Fatalf("AesDecrypt failed,err:%v", err)
			}
			t.Logf("plainText:%s", decryptCode)
		})
	}
}
