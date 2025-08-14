package main

import (
	"exif-extractor/markers"
	"bytes"
	"fmt"
	"log"
	"os"
	"errors"
)

func main() {
//	testingBytes()
//	getExifData()

	getImgData("test-image.jpeg")
}

func getImgData(imgPath string) {
	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatalf("error: %v, is NOT a valid image path", imgPath)
	}
	defer file.Close()

	chunkSize := 4096
	chunk := make([]byte, chunkSize)
	_, err = file.Read(chunk)
	if err != nil {
		log.Fatal("error reading chunk")
	}

	chunkReader := bytes.NewReader(chunk[:])

 	checkMarker(chunkReader, markers.SOIMarker)
 	checkMarker(chunkReader, markers.App0Marker)
	// TODO: check data
	checkMarker(chunkReader, markers.App1Marker)
	// TODO: check data

	dataPayload := make([]byte, 2)
	_, err = chunkReader.Read(dataPayload)
	if err != nil {
		log.Fatal("Unable to read app0 data pay load")
	}
}

func checkMarker(chunk *bytes.Reader, m markers.Marker) error {
	b := make([]byte, 2)
	_, err := chunk.Read(b)
	if err != nil {
		return errors.New("could not read chunk")
	}

	if b[0] != m.First || b[1] != m.Second {
		return fmt.Errorf("Marker: %v not found", m.Name)
	}
	fmt.Printf("Marker: %v found\n", m.Name)

	return nil
}
