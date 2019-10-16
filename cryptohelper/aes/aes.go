package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const (
	_ = iota
	CBC
	CFB
)

type Aes struct {
	mode int8
	key  []byte
	iv   []byte
}

func New(mode int8, key, iv []byte) *Aes {
	return &Aes{mode, key, iv}
}

func (this *Aes) Encrypt(src []byte) ([]byte, error) {
	block, err := aes.NewCipher(this.key)
	if err != nil {
		return nil, err
	}

	b := this.PKCSPadding(src, block.BlockSize())
	data := make([]byte, len(b))

	switch this.mode {
	case CFB:
		bm := cipher.NewCFBEncrypter(block, this.iv)
		bm.XORKeyStream(data, b)
	default:
		bm := cipher.NewCBCEncrypter(block, this.iv)
		bm.CryptBlocks(data, b)
	}

	return data, nil
}

func (this *Aes) Decrypt(src []byte) ([]byte, error) {
	block, err := aes.NewCipher(this.key)
	if err != nil {
		return nil, err
	}

	data := make([]byte, len(src))

	switch this.mode {
	case CFB:
		bm := cipher.NewCFBDecrypter(block, this.iv)
		bm.XORKeyStream(data, src)
	default:
		bm := cipher.NewCBCDecrypter(block, this.iv)
		bm.CryptBlocks(data, src)
	}

	return this.PKCSUnPadding(data), nil
}

// 转换为base64编码后的字符串
func (this *Aes) EncryptToString(src []byte) (string, error) {
	b, err := this.Encrypt(src)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// 解码base64字符串
func (this *Aes) DecryptFromString(base64EncodeStr string) ([]byte, error) {
	b1, err := base64.StdEncoding.DecodeString(base64EncodeStr)
	if err != nil {
		return nil, err
	}

	return this.Decrypt(b1)
}

func (this *Aes) PKCSPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	text := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, text...)
}

func (this *Aes) PKCSUnPadding(src []byte) []byte {
	length := len(src)
	return src[:(length - int(src[length-1]))]
}
