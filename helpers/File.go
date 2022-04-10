package helpers

import (
	"log"
	"os"
)

func GetAllFiles() (os.FileInfo, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	log.Printf(path)
	return nil, nil
}
