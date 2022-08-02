package main

import (
	"fmt"
	"io"
	"log"
	"os"

	aes_cbc "github.com/binsabit/cipher-project/internal/enc_machines/aes-cbc"
	encrypt "github.com/binsabit/cipher-project/internal/encryption"
	"github.com/binsabit/cipher-project/pkg/helpers"

	"github.com/tcolgate/mp3"
)

const (
	helpCmd     = "h"
	encryptCmd  = "e"
	decryptCmd  = "d"
	sendDataCmd = "s"
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
			encm1, err := aes_cbc.New_AES_CBC_Encryption(filePath, secret)
			if err != nil {
				log.Fatal(err)
			}
			_, err = encrypt.EncryptAndSave(encm1, "./r")
			if err != nil {
				log.Fatal(err)
			}
		case cmd == decryptCmd:
			fmt.Println(decryptMsg)
			filePath := helpers.FilePath()
			log.Printf("Started reading file with: %s", filePath)
			decm1, err := aes_cbc.New_AES_CBC_Encryption(filePath, secret)
			if err != nil {
				log.Fatal(err)
			}
			_, err = encrypt.DecryptAndSave(decm1, "./l")
			if err != nil {
				log.Fatal(err)

			}

		default:
			continue
		}
	}
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
