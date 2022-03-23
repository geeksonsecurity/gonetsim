package utils

import (
	"log"
	"os"
)

func WriteBinaryFile(path string, content []byte) {
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bytesWritten, err := file.Write(content)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes to %s.\n", bytesWritten, path)
}
