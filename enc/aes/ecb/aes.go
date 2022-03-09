package ecb

import (
	"crypto/aes"
)

func AesEncrypt(plainText []byte, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return nil, err
	}
	length := (len(plainText) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, plainText)
	pad := byte(len(plain) - len(plainText))
	for i := len(plainText); i < len(plain); i++ {
		plain[i] = pad
	}
	cipherText := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(plainText); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(cipherText[bs:be], plain[bs:be])
	}

	return cipherText, nil
}

func AesDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return nil, err
	}
	plainText := make([]byte, len(cipherText))
	for bs, be := 0, cipher.BlockSize(); bs < len(cipherText); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(plainText[bs:be], cipherText[bs:be])
	}

	trim := 0
	if len(plainText) > 0 {
		trim = len(plainText) - int(plainText[len(plainText)-1])
	}

	return plainText[:trim], nil
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
