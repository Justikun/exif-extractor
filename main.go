package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"github.com/justikun/exif-extractor/image_data"
)

func main() {
	//photos := [3]string{"test-photos/test-image-1.jpeg", "test-photos/test-image-2.tiff", "test-photos/test-image-3"}

	imagesData := []image_data.ImageData{}

	img, err := parseImgData("test-photos/test-image-1.jpeg")
	if err != nil {
		fmt.Println(err)
	}
	imagesData[0] = img
	print(imagesData[0].ImagePath)
}

func parseImgData(imgPath string) (image_data.ImageData, error) {
	imgData := image_data.ImageData{}

	file, err := os.Open(imgPath)
	if err != nil {
		return imgData, fmt.Errorf("error: %v, is NOT a valid image path\n", imgPath)
	}
	defer file.Close()

	for {
		// find marker
		marker := make([]byte, 2)
		_, err = file.Read(marker)
		if err != nil {
			return imgData, err
		}

		if marker[0] != 0xFF {
			return imgData, fmt.Errorf("Hex: %x. Not a valid marker\n", marker[0])
		}

		switch marker[1] {
		case 0xD8: // SOI - start of image
			println("case soi:")
		case 0xD9: // EOI - end of image
			return imgData, fmt.Errorf("case eoi")
		case 0xE0: // APP0 - jfif marker
			println("case app0:")
			parseAPP0(file, &imgData)
		case 0xE1: // APP1
			println("case app1:")
			err := parseAPP1(file, &imgData)
			if err != nil {
				log.Fatal(err)
			}
			println("APP1 DONE")
			println("Data: ", imgData.MetaData.Tags)
			println("--------------")

			return imgData, nil
		case 0xDA: // SOS - image stream
			return imgData, nil
		default:
			return imgData, fmt.Errorf("No marker found")
		}
	}
}

// TODO Find a way to store the data.
func parseAPP0(file *os.File, imgData *image_data.ImageData) error {
	// get data length
	dataLengthB := make([]byte, 2)
	_, err := file.Read(dataLengthB)
	if err != nil {
		return fmt.Errorf("Failed to read data length")
	}
	dataLength := int(binary.BigEndian.Uint16(dataLengthB)) - 2

	// update main file reader
	buf := make([]byte, dataLength)
	_, err = file.Read(buf)

	// create new reader
	segmentReader := bytes.NewReader(buf)

	// check JFIF identifier
	identifier := make([]byte, 5)
	_, err = segmentReader.Read(identifier)
	if err != nil {
		return err
	}
	if string(identifier) != "JFIF\x00" {
		log.Fatal("Not JFIF compliant")
	}

	// get version
	version, err := segmentReader.Read(make([]byte, 2))
	if err != nil {
		return fmt.Errorf("APP0: Failed to read version")
	}
	units, err  := segmentReader.Read(make([]byte, 1)) // 00 - no units / 01 - DPI / 02 - Dots per centimeter
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
	// I have pre-read this data with the main file reader in order to continue

	// TODO: Read non standard data
	// TOOD: Save Data to an objetc
	//
	return nil
}


func parseAPP1(file *os.File, imgData *image_data.ImageData) error {
	// read payload size
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

	// check version - TIFF magic number 42
	versionNumber := endian.Uint16(tiffHeader[2:4])
	if versionNumber != 42 {
		return errors.New("invalid version number")
	}
	// update file reader to offset of the first IFD (Image File Directory)
	ifdOffset := endian.Uint32(tiffHeader[4:8])
	if _, err = file.Seek(int64(ifdOffset), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to offset %d: %v", ifdOffset, err)
	}
	err = parseIFD(app1Reader, tiffHeaderStart, image_data.IFDMAIN, endian)
	if err != nil {
		return err
	}
	return nil
}

func parseIFD(reader *bytes.Reader, tiffHeaderStart int64, ifdType image_data.IFDType, endian binary.ByteOrder) error {
	// count of tags
	tagInBytes:= make([]byte, 2)
	_, err := reader.Read(tagInBytes)
	if err != nil {
		return err
	}
	tagCount := int(endian.Uint16(tagInBytes))

	for i := 0; i < tagCount; i++ {

		tag := image_data.IFDTag{}

		// set id
		idInBytes := make([]byte, 2)
		_, err = reader.Read(idInBytes)
		if err != nil {
			return err
		}
		tag.Id = endian.Uint16(idInBytes)

		// set tag name
		tagName, err := image_data.GetNameFromID(endian.Uint16(idInBytes))
		tag.Name = tagName

		// set data type
		dataTypeB := make([]byte, 2)
		_, err = reader.Read(dataTypeB)
		if err != nil {
			return fmt.Errorf("Failed to read dataTypeB\n")
		}

		dataType, err := image_data.GetDataTypeFromBytes(dataTypeB, endian)
		if err != nil {
			return err
		}
		tag.DataType = dataType

		// set count of data
		countDataInBytes := make([]byte, 4)
		_, err = reader.Read(countDataInBytes)
		if err != nil {
			return fmt.Errorf("Failed to read count buff\n")
		}
		tagCount := endian.Uint32(countDataInBytes)
		tag.Count = tagCount

		// check data size
		dataTypeSize, err := dataType.SizeInBytes()
		if err != nil {
			return err
		}
		totalTagSize := dataTypeSize * int(tagCount)

		if totalTagSize > 4 {
			// next 4 bytes is a pointer to the data
			//tag.Data = getDataFromPointer(offset)
			fmt.Println("DATA OUTSIDE OF TAG DATA")
			dataOffset := make([]byte, 4)
			_, err = reader.Read(dataOffset)
			if err != nil {
				return fmt.Errorf("Error reading dataOffset")

			}
			fmt.Printf("dataOffset: %x\n", dataOffset)
		} else {
			// next 4 bytes holds data
			fmt.Println("DATA INSIDE OF TAG DATA")
			data := make([]byte, 4)
			_, err = reader.Read(data)
			if err != nil {
				return err
			}
			tag.Data = data
		}

		///////// DEBUG PRINTING ////////////
		dataTypeString, err := tag.DataType.GetDataTypeString()
		if err != nil {
			return err
		}

		fmt.Printf("id: %v, name: %v, dataType: %v, dataValue: %s", tag.Id, tag.Name, dataTypeString, tag.Data)
		/////////END OF DEBUG////////////

		// parse and set value
	}
	print("END OF PARSEID")
	return nil
}
