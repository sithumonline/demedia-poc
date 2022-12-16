package utility

import (
	"fmt"
	"github.com/sithumonline/demedia-poc/core/config"
	"log"
	"os"
)

func WriteFile(data string) {
	err := os.WriteFile(config.AddressFilePath, []byte(fmt.Sprintf("%s", data)), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func ReadFile() string {
	dat, err := os.ReadFile(config.AddressFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	return string(dat)
}
