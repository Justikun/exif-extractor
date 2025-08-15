package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	photos := [3]string{"test-photos/test-image-1.jpeg", "test-photos/test-image-2.tiff", "test-photos/test-image-3"}

	err := getImgData(photos[0])
	if err != nil {
		fmt.Println(err)
	}
}

func getImgData(imgPath string) error {
	file, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("error: %v, is NOT a valid image path\n", imgPath)
	}
	defer file.Close()

	chunkSize := 4096
	chunk := make([]byte, chunkSize)
	_, err = file.Read(chunk)
	if err != nil {
		return fmt.Errorf("error reading chunk into 4096\n")
	}

	chunkReader := bytes.NewReader(chunk[:])

	// Loop over chunkData until EOI
	for {
		// Read marker
		marker := make([]byte, 2)
		_, err := chunkReader.Read(marker)
		if err != nil {
			return fmt.Errorf("couln't read 2 byte chunk in chunkReader\n")
		}
		if marker[0] != 0xFF {
			//print("%x: not a valid marker\n", marker[0])
			log.Fatalf("Hex: %x. Not a valid marker\n", marker[0])
		}

		switch marker[1] {
		case 0xD8: // SOI - start of image
		fmt.Printf("found SOI\n")
		case 0xD9: // EOI - end of image
		fmt.Printf("found EOI\n")
		case 0xE0: // APP0 - jfif marker
		fmt.Printf("found APPO\n")
		if err := parseAPP0(chunkReader); err != nil {
			log.Fatalln(err)
		}
		case 0xE1: // APP1
		fmt.Printf("found APP1\n")
		if err := parseAPP1(chunkReader); err != nil {
			log.Fatal(err)
		}
		case 0xDA: // SOS - image stream
		}
	}
}

// TODO Find a way to store the data.
func parseAPP0(r *bytes.Reader) error {
	// check payload size
	b := make([]byte, 2)
	_, err := r.Read(b)
	if err != nil {
		return fmt.Errorf("coudln't check payload size\n")
	}

	// Go has a std lib for this
	dataSizeInBytes := int(b[0])<<8 | int(b[1])
	data := make([]byte, dataSizeInBytes-2)
	_, err = r.Read(data)
	if err != nil {
		return err
	}

	var dataString string = byteSliceToString(data)
	fmt.Printf("data: %v\n", dataString)
	return nil
}


func parseAPP1(r *bytes.Reader) error {
	b := make([]byte, 2)
	_, err := r.Read(b)
	if err != nil {
		return err
	}

	// TODO Look into APP1. What markers can exists.
	// TODO Find a way to store the data
	// TODO I'm assuming the next two bytes indicate how large app1 segment is.Double check this
	dataSizeInBytes := binary.BigEndian.Uint16(b[:])
	data := make([]byte, dataSizeInBytes)
	_, err = r.Read(data)
	if err != nil {
		return err
	}

	var dataString string = byteSliceToString(data)
	fmt.Printf("data: %v\n", dataString)
	return nil
}

func byteSliceToString(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	results := make([]string, len(data))
	for i, b := range data {
		results[i] = fmt.Sprintf("%c", b)
	}

	return strings.Join(results, "")
}
