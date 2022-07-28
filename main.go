package main

import (
	"crypto/aes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/binsabit/cipher-project/internal/data"
	"github.com/binsabit/cipher-project/internal/encryption"
	"github.com/binsabit/cipher-project/pkg/helpers"

	"github.com/tcolgate/mp3"
)

//IV + json size + json + file

const (
	helpCmd    = "h"
	encryptCmd = "e"
	decryptCmd = "d"
)

const (
	helpMsg    = "Please choose command: h -help , e -encrypt, d -decrypt"
	encryptMsg = "[ENCRYPTION MODE]:\nPlease enter path to file"
	decryptMsg = "[DECRYPTION MODE]:\nPlease enter path to file"
)

func main() {
	secret := helpers.EnvVars()["SECRET"]

	for {
		fmt.Println(helpMsg)
		var cmd string
		fmt.Scanln(&cmd)
		switch {
		case cmd == helpCmd:
			continue
		case cmd == encryptCmd:
			fmt.Println(encryptMsg)
			filePath := helpers.FilePath()
			log.Printf("Started reading file with: %s", filePath)

			file, err := helpers.ReadFile(filePath)
			if err != nil {
				log.Println(err)
				continue
			}
			metadata, err := helpers.ToJson(data.GetMetadata(filePath))
			if err != nil {
				log.Println("could not process the data: json marshal")
				continue
			}

			encrptedData, IV, err := encryption.Encrypt(file, secret)
			if err != nil {
				fmt.Println("Error Ecryption")

			}
			fmt.Println(IV)
			metaSize := make([]byte, 4)
			binary.LittleEndian.PutUint32(metaSize, uint32(len(metadata)))
			fmt.Println(metaSize)
			var a data.Metadata
			err = json.Unmarshal(metadata, &a)
			if err != nil {
				log.Println("could not process the data: json unmarshal")
			}
			fmt.Println(a)
			h := sha1.Sum(encrptedData)
			hash := h[:]
			res := addAll(IV, metaSize, metadata, hash, encrptedData)
			_ = os.Mkdir("./test-files-result", os.ModePerm)
			os.Create("./test-files-result/" + a.Name)
			ioutil.WriteFile("./test-files-result/"+a.Name, res, os.ModePerm)
		case cmd == decryptCmd:

			fmt.Println(decryptMsg)
			filePath := helpers.FilePath()
			log.Printf("Started reading file with: %s", filePath)

			file, err := helpers.ReadFile(filePath)
			if err != nil {
				log.Println(err)
				continue
			}
			IV := file[:aes.BlockSize]
			jsonSizeByte := file[aes.BlockSize : aes.BlockSize+4]
			jsonSizeInt := int(binary.LittleEndian.Uint32(jsonSizeByte))
			metadata := file[aes.BlockSize+4 : aes.BlockSize+4+jsonSizeInt]
			var a data.Metadata
			err = json.Unmarshal(metadata, &a)
			if err != nil {
				log.Println("could not process the data: json unmarshal")
			}
			fmt.Println(a)
			hash := file[aes.BlockSize+4+jsonSizeInt : aes.BlockSize+4+jsonSizeInt+20]
			content := file[aes.BlockSize+4+jsonSizeInt+20:]
			h := sha1.Sum(content)
			if string(hash) != string(h[:]) {
				log.Fatal("Data lost")
			}
			decryptRes, err := encryption.Decrypt(content, []byte(secret), IV)
			if err != nil {
				log.Println("could not decrypt")
				continue
			}
			_ = os.Mkdir("./result", os.ModePerm)
			os.Create("./result/" + a.Name)
			ioutil.WriteFile("./result/"+a.Name, decryptRes, os.ModePerm)

		default:
			continue
		}
	}

	// log.Println(file)
}
func addAll(IV, metasize, metadata, hash, encrptedData []byte) []byte {
	var res []byte
	res = append(res, IV...)
	res = append(res, metasize...)
	res = append(res, metadata...)
	res = append(res, hash...)
	res = append(res, encrptedData...)
	return res
}

func Duration(mp3File string) {
	t := 0.0

	r, err := os.Open(mp3File)
	if err != nil {
		fmt.Println(err)
		return
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {

		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}

		t = t + f.Duration().Seconds()
	}
	fmt.Println(t)
}
