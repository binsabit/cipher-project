package aescbc

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/binsabit/cipher-project/internal/data"
	"github.com/binsabit/cipher-project/internal/encryption"
	"github.com/binsabit/cipher-project/pkg/helpers"
)

const (
	jsonMemBytes = 4
	hashMemBytes = 20
)

type AES_CBC_Encryption struct {
	content  []byte        //content of file
	metadata data.Metadata //json encodedd metadata of file
	key      string        //secret key for encryption
}

func New_AES_CBC_Encryption(filepath, key string) (*AES_CBC_Encryption, error) {
	content, err := helpers.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	metadata := data.GetMetadataOfFile(filepath)
	return &AES_CBC_Encryption{
		content:  content,
		metadata: metadata,
		key:      key,
	}, nil
}

func (a AES_CBC_Encryption) Encrypt() (encryptedContent []byte, err error) {
	hash := helpers.HashData(a.content)

	metadataJson, err := helpers.ToJson(a.metadata)
	if err != nil {
		return nil, err
	}
	metaSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(metaSize, uint32(len(metadataJson)))
	log.Println(metaSize)
	IV := encryption.GenerateIV(aes.BlockSize)
	fullcontent := addAll(metaSize, metadataJson, hash, a.content)
	log.Println(IV)
	log.Println(hash)

	bPlaintext := encryption.PKCS5Padding(fullcontent, aes.BlockSize)
	block, err := aes.NewCipher([]byte(a.key))

	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(bPlaintext))

	mode := cipher.NewCBCEncrypter(block, IV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return append(IV, ciphertext...), nil
}

func (a AES_CBC_Encryption) EncryptAndSave(destination string) (encryptedContent []byte, err error) {
	data, err := a.Encrypt()
	if err != nil {
		return nil, err
	}
	_ = helpers.CreateFileAndWrite(destination, a.metadata.Name, data)

	return data, nil
}

func (a AES_CBC_Encryption) Decrypt() (decryptedContent []byte, err error) {
	block, err := aes.NewCipher([]byte(a.key))

	if err != nil {
		return nil, err
	}
	IV := encryption.IVofFile(a.content)
	log.Println(IV)
	mode := cipher.NewCBCDecrypter(block, IV)
	dec := make([]byte, len(a.content))

	mode.CryptBlocks(dec, a.content)
	// cutTrailingSpaces := []byte(strings.TrimSpace(string(dec)))

	jsonSizeByte := dec[aes.BlockSize : aes.BlockSize+jsonMemBytes]
	log.Println(jsonSizeByte)
	jsonSizeInt := int(binary.LittleEndian.Uint32(jsonSizeByte))

	metadata := dec[aes.BlockSize+jsonMemBytes : aes.BlockSize+jsonMemBytes+jsonSizeInt]

	data, err := data.GetMetadataFromJson(metadata)
	if err != nil {
		return nil, err
	}
	log.Println(data)
	hashPrev := dec[aes.BlockSize+jsonMemBytes+jsonSizeInt : aes.BlockSize+jsonMemBytes+jsonSizeInt+hashMemBytes]
	log.Println(hashPrev)

	content := []byte(strings.TrimSpace(string(dec[aes.BlockSize+jsonSizeInt+jsonMemBytes+hashMemBytes : aes.BlockSize+jsonSizeInt+jsonMemBytes+hashMemBytes+int(data.Size)])))

	hashCur := helpers.HashData(content)

	if !helpers.MatchHash(hashPrev, hashCur) {
		return nil, fmt.Errorf("file damaged hashes does not match")
	}

	return content, err
}

func (a AES_CBC_Encryption) DecryptAndSave(destination string) ([]byte, error) {
	data, err := a.Decrypt()
	if err != nil {
		return nil, err
	}
	_ = helpers.CreateFileAndWrite(destination, a.metadata.Name, data)
	return data, nil
}

func addAll(metasize, metadata, hash, data []byte) []byte {
	var res []byte
	res = append(res, metasize...)
	res = append(res, metadata...)
	res = append(res, hash...)
	res = append(res, data...)
	return res
}
