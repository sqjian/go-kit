package sha1

import (
	"crypto/sha1"
	"encoding/hex"
)

func Calc(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
