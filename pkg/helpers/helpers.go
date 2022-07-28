package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
