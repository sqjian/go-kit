package aes

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
		{name: "test1", args: args{plainText: []byte("hello world"), key: []byte("1443flfsaWfdasds")}},
	}

	aes, _ := NewAes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("plainText：%s,key:%s", tt.args.plainText, tt.args.key)
			cipherText, err := aes.AesEncrypt(tt.args.plainText, tt.args.key)
			if err != nil {
				t.Fatalf("AesEncrypt failed,err:%v", err)
			}
			t.Logf("cipherText data:%X", cipherText)
			plainText, err := aes.AesDecrypt(cipherText, tt.args.key)
			if err != nil {
				t.Fatalf("AesDecrypt failed,err:%v", err)
			}
			t.Logf("plainText:%s", plainText)
		})
	}
}
