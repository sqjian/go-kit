package aes

import (
	"fmt"
	"github.com/sqjian/go-kit/crypto/aes/cbc"
	"github.com/sqjian/go-kit/crypto/aes/cfb"
	"github.com/sqjian/go-kit/crypto/aes/ctr"
	"github.com/sqjian/go-kit/crypto/aes/ecb"
	"github.com/sqjian/go-kit/crypto/aes/ofb"
)

type AES interface {
	AesEncrypt([]byte, []byte) ([]byte, error)
	AesDecrypt([]byte, []byte) ([]byte, error)
}

type aes struct {
	mode Mode
}

func (a aes) AesEncrypt(plainText []byte, key []byte) ([]byte, error) {
	switch a.mode {
	case ECB:
		return ecb.AesEncrypt(plainText, key)
	case CBC:
		return cbc.AesEncrypt(plainText, key)
	case CFB:
		return cfb.AesEncrypt(plainText, key)
	case CTR:
		return ctr.AesEncrypt(plainText, key)
	case OFB:
		return ofb.AesEncrypt(plainText, key)
	default:
		return nil, fmt.Errorf("internal error,unknown mode:%v", a.mode.String())
	}
}

func (a aes) AesDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	switch a.mode {
	case ECB:
		return ecb.AesDecrypt(cipherText, key)
	case CBC:
		return cbc.AesDecrypt(cipherText, key)
	case CFB:
		return cfb.AesDecrypt(cipherText, key)
	case CTR:
		return ctr.AesDecrypt(cipherText, key)
	case OFB:
		return ofb.AesDecrypt(cipherText, key)
	default:
		return nil, fmt.Errorf("internal error,unknown mode:%v", a.mode.String())
	}
}

func newDefaultAesConfig() *aes {
	return &aes{mode: ECB}
}

func NewAes(opts ...Option) (AES, error) {

	aes := newDefaultAesConfig()

	for _, opt := range opts {
		opt.apply(aes)
	}
	return aes, nil
}
