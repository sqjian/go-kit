package ofb

import (
	"bytes"
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
	plainText = PKCS7Padding(plainText, blockSize)
	cipherText := make([]byte, blockSize+len(plainText))
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(cipherText[blockSize:], plainText)
	return cipherText, nil
}

func AesDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	iv := cipherText[:blockSize]
	cipherText = cipherText[blockSize:]
	if len(cipherText)%blockSize != 0 {
		return nil, fmt.Errorf("cipherText is not a multiple of the block size")
	}

	plainText := make([]byte, len(cipherText))
	mode := cipher.NewOFB(block, iv)
	mode.XORKeyStream(plainText, cipherText)

	plainText = PKCS7UnPadding(plainText)
	return plainText, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	upPadding := int(origData[length-1])
	return origData[:(length - upPadding)]
}
