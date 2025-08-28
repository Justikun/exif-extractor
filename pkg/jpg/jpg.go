package jpg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/justikun/metadata-viewer/pkg/metadata"
	"github.com/justikun/metadata-viewer/pkg/tiff"
)

func ParseAPP0(file *os.File, imgData *metadata.ImageData) error {
	// get data length
	dataLengthB := make([]byte, 2)
	_, err := file.Read(dataLengthB)
	if err != nil {
		return fmt.Errorf("Failed to read data length")
	}
	dataLength := int(binary.BigEndian.Uint16(dataLengthB)) - 2

	// update main file app1Reader
	buf := make([]byte, dataLength)
	_, err = file.Read(buf)

	// create new app1Reader
	segmentReader := bytes.NewReader(buf)

	// check JFIF identifier
	identifier := make([]byte, 5)
	_, err = segmentReader.Read(identifier)
	if err != nil {
		return err
	}
	if string(identifier) != "JFIF\x00" {
		return fmt.Errorf("Not JFIF compliant")
	}

	// get version
	version, err := segmentReader.Read(make([]byte, 2))
	if err != nil {
		return fmt.Errorf("APP0: Failed to read version")
	}
	units, err := segmentReader.Read(make([]byte, 1)) // 00 - no units / 01 - DPI / 02 - Dots per centimeter
	if err != nil {
		return fmt.Errorf("APP0: Failed to read units")
	}
	xDensity, err := segmentReader.Read(make([]byte, 2)) // horizontal pixel density
	if err != nil {
		return fmt.Errorf("APP0: Failed to read xDensity")
	}
	yDensity, err := segmentReader.Read(make([]byte, 2)) // vertical pixel density
	if err != nil {
		return fmt.Errorf("APP0: Failed to read yDensity")
	}
	thumbnW, err := segmentReader.Read(make([]byte, 1)) // 00 is no thumbnail
	if err != nil {
		return fmt.Errorf("APP0: Failed to read thumbnail width")
	}
	thumbnH, err := segmentReader.Read(make([]byte, 1)) // 00 is no thumbnail
	if err != nil {
		return fmt.Errorf("APP0: Failed to read thumbnail height")
	}
	// w * h * 3 = x bytes
	thumbnBytes := thumbnW * thumbnH * 3
	thumbnData, err := segmentReader.Read(make([]byte, thumbnBytes))
	if err != nil {
		return fmt.Errorf("APP0: Failed to read thumbnail data")
	}

	fmt.Printf("\n\n\nSpitting Data:%v%v%v%v%v%v%v%v\n\n\n", thumbnData, version, units, xDensity, yDensity, thumbnW, thumbnH, identifier)
	// There can be other non standard data past this part
	// I have not parsed it yet
	// I have pre-read this data with the main file app1Reader in order to continue

	// TODO: Read non standard data
	// TOOD: Save Data to an objetc
	//
	return nil
}

func ParseAPP1(file *os.File, imgData *metadata.ImageData) error {
	// check APP1 payload size
	buff := make([]byte, 2)
	_, err := io.ReadFull(file, buff)
	if err != nil {
		return err
	}

	// read entire APP1 and advance main file reader
	payloadSize := binary.BigEndian.Uint16(buff)
	app1Data := make([]byte, payloadSize-2)
	file.Read(app1Data)

	app1Reader := bytes.NewReader(app1Data)

	// check identifier
	identifier := make([]byte, 6)
	_, err = app1Reader.Read(identifier)
	if err != nil {
		return err
	}
	payloadIdentifier := string(identifier)

	var endian binary.ByteOrder
	switch payloadIdentifier {
	case "Exif\x00\x00":
		endian = binary.BigEndian
	case "http://ns.adobe.com/xap/1.0/\x00":
		endian = binary.LittleEndian
	default:
		return errors.New("Payload identifier not found")
	}

	tiffHeaderStart, err := app1Reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("error setting tiff header start")
	}
	fmt.Printf("tiff header in app1 starts at: %d\n", tiffHeaderStart)
	// read tiff header
	tiffHeader := make([]byte, 8)
	_, err = app1Reader.Read(tiffHeader)
	if err != nil {
		return fmt.Errorf("Error reading tiff header")
	}

	// check endianness (byte order)
	endianString := string(tiffHeader[:2])
	switch endianString {
	case "MM":
		endian = binary.BigEndian
	case "II":
		endian = binary.LittleEndian
	default:
		return errors.New("Endianess not found")
	}

	// wrap app1Reader into BinaryReader
	br := metadata.NewBinaryReader(app1Reader, endian)

	// check version - TIFF magic number 42
	versionNumber := endian.Uint16(tiffHeader[2:4])
	if versionNumber != 42 {
		return errors.New("invalid version number")
	}
	// update file main reader to offset of the first IFD (Image File Directory)
	ifdOffset := endian.Uint32(tiffHeader[4:8])
	if _, err = file.Seek(int64(ifdOffset), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to offset %d: %v", ifdOffset, err)
	}
	err = tiff.ParseIFD(br, tiffHeaderStart, metadata.IFDMAIN, endian)
	if err != nil {
		return err
	}
	return nil
}
