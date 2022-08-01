package helpers

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	ErrCantReadFile = errors.New("cannot read file with thod path")
)

func ToJson(item interface{}) ([]byte, error) {
	res, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func EnvVars() map[string]string {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("could not read .env file", err)
	}
	return envMap
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file:", err)
	}
	return file, nil
}

func FilePath() string {
	var path string
	fmt.Scanln(&path)
	return path
}

func HashData(data []byte) []byte {
	hash := sha1.Sum(data)
	return hash[:]
}

func MatchHash(h1, h2 []byte) bool {
	return string(h1) == string(h2)
}

func CreateFileAndWrite(destination, filename string, data []byte) error {
	_ = os.Mkdir(destination, os.ModePerm)
	fPath := filepath.Join(destination, filename)
	os.Create(fPath)
	ioutil.WriteFile(fPath, data, os.ModePerm)
	return nil
}
