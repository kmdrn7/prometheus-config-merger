package utils

import (
	"log"
	"os"
)

func SyncResourceContentToLocalFile(content []byte, filepath string) error {
	log.Printf("overwriting file [%s] with new content \n", filepath)
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		panic(err)
	}
	return nil
}
