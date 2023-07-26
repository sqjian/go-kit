package simplifiedchinese

import (
	"testing"
)

func TestConvert(t *testing.T) {
	utf8Byte := []byte("UTF-8字符包子")
	t.Log(string(utf8Byte))

	// UTF-8 转 BIG5
	s, _ := EncodeBig5(utf8Byte)
	t.Log(string(s))

	// BIG5 转 UTF-8
	s, _ = DecodeBig5(s)
	t.Log(string(s))

	// UTF-8 转 GBK
	s, _ = EncodeGBK(s)
	t.Log(string(s))

	// GBK 转 UTF-8
	s, _ = DecodeGBK(s)
	t.Log(string(s))
}
