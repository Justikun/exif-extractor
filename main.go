package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
//	testingBytes()
//	getExifData()


	getImgData("test-image.jpeg")
}


func getImgData(imgPath string) {
	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatal("error: %v, is NOT a valid image path", imgPath)
	}
	defer file.Close()

	chunkSize := 4096
	chunk := make([]byte, chunkSize)
	_, err = file.Read(chunk)
	if err != nil {
		log.Fatal("error reading chunk")
	}

	reader := bytes.NewReader(chunk[:])

	// check 0xFFD8 - JPG - SOI (start of image)
	header := make([]byte, 2)
	_, err = reader.Read(header)
	if err != nil {
		log.Fatal("error reading header")
	}
	if header[0] != 0xFF && header[1] != 0xD8 {
		log.Fatal("file is not JPG")
	}
	fmt.Printf("%x\n", header[0])
	fmt.Printf("%x\n", header[1])
}
