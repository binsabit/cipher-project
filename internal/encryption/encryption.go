package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"strings"
)

func Encrypt(content []byte, key string) ([]byte, []byte, error) {
	bPlaintext := pkCS5Padding(content, aes.BlockSize)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, nil, err
	}

	IV := generateIV(aes.BlockSize)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, IV)
	mode.CryptBlocks(ciphertext, bPlaintext)

	return ciphertext, IV, nil
}

func Decrypt(cipherText []byte, encryptionKey []byte, IV []byte) (decryptedContent []byte, err error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, IV)
	mode.CryptBlocks(cipherText, cipherText)

	cutTrailingSpaces := []byte(strings.TrimSpace(string(cipherText)))
	return cutTrailingSpaces, err
}

func pkCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(cipherText, padText...)
}

func generateIV(size int) []byte {
	iv := make([]byte, size)
	rand.Read(iv)
	return iv
}
