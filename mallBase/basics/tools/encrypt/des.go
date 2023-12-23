package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

func padding(src []byte, blocksize int) []byte {
	n := len(src)
	padnum := blocksize - n%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	dst := append(src, pad...)
	return dst
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	dst := src[:n-unpadnum]
	return dst
}

func EncryptDES(src []byte) string {
	key := []byte("12345678")
	block, _ := des.NewCipher(key)
	src = padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return base64.StdEncoding.EncodeToString(src)
}

func DecryptDES(src string) []byte {
	key := []byte("12345678")
	block, _ := des.NewCipher(key)
	secretData, _ := base64.StdEncoding.DecodeString(src)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(secretData, secretData)
	return unpadding(secretData)
}
