package ctr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func AesEncrypt(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	cipherText := make([]byte, blockSize+len(plainText))
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[blockSize:], plainText)
	return cipherText, nil
}

func AesDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := cipherText[:blockSize]
	cipherText = cipherText[blockSize:]

	plainText := make([]byte, len(cipherText))

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plainText, cipherText)

	return plainText, nil
}
