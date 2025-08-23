package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"github.com/justikun/metadata-viewer/pkg/metadata"
	"path/filepath"
	"strings"
)

func main() {
	//photos := [3]string{"test-photos/test-image-1.jpeg", "test-photos/test-image-2.tiff", "test-photos/test-image-3"}

	images, err := GetImageFiles("test-photos")
	if err != nil {
		fmt.Println(err)
	}

	for _, image := range images {
		images, err := parseImgData(image.ImagePath)
		if err != nil {
			fmt.Printf("Failed to parse image at %s", images.ImagePath)
		}
	}

	print(images[0].ImagePath)
}

func GetImageFiles(dirPath string) ([]metadata.ImageData, error) {
	var imageFiles []metadata.ImageData

	allowedExtensions := map[string]struct{} {
		".jpg": {},
		".jpeg": {},
		".png": {},
		".tiff": {},
		".tif": {},
		".webp": {},
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return []metadata.ImageData{}, fmt.Errorf("Failed to read directory %s: %w", dirPath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		_, ok := allowedExtensions[ext]
		if ok {
			// construct absolute path
			absPath := filepath.Join(dirPath, fileName)
			image := metadata.ImageData{ImagePath: absPath, MetaData: metadata.MetaData{}}
	 		imageFiles = append(imageFiles, image)
		}
	}
	return imageFiles, nil
}

func parseImgData(imgPath string) (metadata.ImageData, error) {
	imgData := metadata.ImageData{}

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
func parseAPP0(file *os.File, imgData *metadata.ImageData) error {
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
	// I have pre-read this data with the main file app1Reader in order to continue

	// TODO: Read non standard data
	// TOOD: Save Data to an objetc
	//
	return nil
}


func parseAPP1(file *os.File, imgData *metadata.ImageData) error {
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
	err = parseIFD(br, tiffHeaderStart, metadata.IFDMAIN, endian)
	if err != nil {
		return err
	}
	return nil
}

func parseIFD(br *metadata.BinaryReader, tiffHeaderStart int64, ifdType metadata.IFDType, endian binary.ByteOrder) error {
	// wrap reader in custom reader
	//TODO
	// count of tags
	tagInBytes, err := br.ReadBytes(2)
	if err != nil {
		return err
	}
	tagCount := int(endian.Uint16(tagInBytes))

	for i := 0; i < tagCount; i++ {
		tag := metadata.IFDTag{}
		// set id
		idInBytes, err := br.ReadBytes(2)
		if err != nil {
			return err
		}
		tag.Id = endian.Uint16(idInBytes)

		// set tag name
		tagName, err := metadata.GetNameFromID(endian.Uint16(idInBytes))
		tag.Name = tagName

		// set data type
		dataTypeB, err := br.ReadBytes(2)
		if err != nil {
			return fmt.Errorf("Failed to read dataTypeB\n")
		}
		dataType, err := metadata.GetDataType(dataTypeB, endian)
		if err != nil {
			return err
		}
		tag.DataType = dataType

		// set count of data
		countDataInBytes, err :=  br.ReadBytes(4)
		if err != nil {
			return fmt.Errorf("Failed to read count of data in bytes\n")
		}
		tagCount := endian.Uint32(countDataInBytes)
		tag.Count = tagCount

		// check data size
		dataTypeSize, err := dataType.SizeInBytes()
		if err != nil {
			return err
		}
		totalTagSize := dataTypeSize * int(tagCount)

		dataOrOffset, err := br.ReadBytes(4)
		if err != nil {
			return fmt.Errorf("Failed to read dataOrOffset")
		}

		if totalTagSize > 4 {
			fmt.Printf("Data > 4 bytes\n")
			// save current pos
			offset := endian.Uint32(dataOrOffset)
			currentPos, err := br.Seek(0, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("Failed to save current position\n")
			}
			// jump to offset
			_, err = br.Seek(int64(offset), io.SeekStart)
			if err != nil {
				return fmt.Errorf("Failed to seek to app1 data tag offset\n")
			}
			// Read the data in bytes
			dataInBytes, err := br.ReadBytes(totalTagSize)
			if err != nil {
				fmt.Printf("total Tag size: %v\n", totalTagSize)
				return fmt.Errorf("Failed to read dataInBytes\n")
			}
			// Parse data
			// TODO
			fmt.Printf("count: %d, size: %d", tagCount, totalTagSize)
			fmt.Printf(" Parsing data at offset: \n")
			fmt.Printf(string(dataInBytes))
			//metadata.DecodeTagByteData(dataInBytes, endian)

			// set reader back to original pos
	 		br.Seek(currentPos, io.SeekStart)
		} else {
			// next 4 bytes holds data
			fmt.Printf("Data < 4 bytes")
			tag.Data = dataOrOffset
		}

		///////// DEBUG PRINTING ////////////
		dataTypeString, err := tag.DataType.GetDataTypeString()
		if err != nil {
			return err
		}

		fmt.Printf("id: %v, name: %v, dataType: %v, dataValue: %s\n\n", tag.Id, tag.Name, dataTypeString, tag.Data)
		/////////END OF DEBUG////////////

		// parse and set value
	}
	print("END OF PARSEID")
	return nil
}
