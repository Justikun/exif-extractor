package main

import (
	"bytes"
	"fmt"
)

func main() {
	testingBytes()
	getExifData()
}

func testingBytes() {
	data := []byte("this is some data")
	fmt.Printf("T: %T\n", data) //T: []uint8
	fmt.Println(data)
	fmt.Printf("\n---------\n\n")

	// byte library - func NewReader(b []byte) *Reader
	reader := bytes.NewReader(data)
	fmt.Printf("T: %T\n", reader) // T: *bytes.Reader
	fmt.Println(reader)
	fmt.Printf("\n---------\n\n")
}

func getExifData() {
	buf := make([]byte, 1000000000)
	fmt.Printf("T: %T\n", buf)
}
