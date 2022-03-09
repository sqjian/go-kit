package ctr

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func AesEncrypt(plainText []byte, key []byte) ([]byte, error) {
	//指定使用AES加密算法
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//给到一个iv值，长度等于block的块尺寸，即16byte
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	//CTR模式是不需要填充的，返回一个计数器模式的，底层采用block生成key流的srtream接口
	stream := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(plainText))
	//加密操作
	stream.XORKeyStream(cipherText, plainText)
	return cipherText, nil
}

func AesDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	//1.指定算法:aes
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//2.返回一个计数器模式的、底层采用block生成key流的Stream接口，初始向量iv的长度必须等于block的块尺寸。
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)

	//3.解密操作
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)

	return plainText, nil
}
