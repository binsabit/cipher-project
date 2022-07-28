package data

import (
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

func GetMetadata(file string) Metadata {
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
