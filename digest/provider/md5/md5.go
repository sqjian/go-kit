package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func Calc(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
