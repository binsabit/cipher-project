package encryption

import (
	"bytes"
	"crypto/aes"
	"math/rand"
)

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(cipherText, padText...)
}

func GenerateIV(size int) []byte {
	iv := make([]byte, size)
	rand.Read(iv)
	return iv
}

func IVofFile(file []byte) []byte {
	return file[:aes.BlockSize]
}
