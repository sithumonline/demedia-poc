package utility

import (
	"errors"
	"fmt"
	"github.com/sithumonline/demedia-poc/core/config"
	"io/fs"
	"log"
	"os"
)

func WriteFile(data string, path string) {
	if path == "" {
		path = config.AddressFilePath
	}
	err := os.WriteFile(path, []byte(fmt.Sprintf("%s", data)), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func ReadFile(path string) string {
	if path == "" {
		path = config.AddressFilePath
	}
	dat, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "file_does_not_exist"
	} else if err != nil {
		log.Fatalln(err)
	}
	return string(dat)
}
