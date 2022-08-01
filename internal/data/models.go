package data

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/h2non/filetype"
)

type Metadata struct {
	Size      int64         `json:"size"`
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	Extension string        `json:"extension"`
	Duration  time.Duration `json:"duration"`
}

func GetMetadataOfFile(file string) Metadata {
	meta, err := os.Stat(file)
	if err != nil {
		log.Fatal("could not read meta data of file:", file, err)
	}
	e := strings.Split(file, ".")
	extension := e[2]
	return Metadata{
		Name:      meta.Name(),
		Type:      filetype.GetType(extension).MIME.Value,
		Extension: extension,
		Size:      meta.Size(),
	}
}

func GetMetadataFromJson(metadata []byte) (Metadata, error) {
	var s Metadata
	err := json.Unmarshal(metadata, &s)
	if err != nil {
		return Metadata{}, err
	}
	return s, nil
}
