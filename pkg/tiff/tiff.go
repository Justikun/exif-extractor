package tiff

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/justikun/metadata-viewer/pkg/metadata"
)

func ParseIFD(br *metadata.BinaryReader, tiffHeaderStart int64, ifdType metadata.IFDtype, endian binary.ByteOrder) error {
	// wrap reader in custom reader
	// TODO

	// count of tags
	tagInBytes, err := br.ReadBytes(2)
	if err != nil {
		return err
	}
	tagCount := int(endian.Uint16(tagInBytes))

	for range tagCount {
		tag := metadata.IFDtag{}

		// set id
		idBytes, err := br.ReadBytes(2)
		if err != nil {
			return err
		}
		tag.ID = endian.Uint16(idBytes)

		// set tag name
		tagName, err := metadata.GetNameFromID(endian.Uint16(idBytes))
		tag.Name = tagName

		// set data type
		dataTypeBytes, err := br.ReadBytes(2)
		if err != nil {
			return fmt.Errorf("Failed to read dataTypeB\n")
		}
		dataType, err := metadata.GetDataType(dataTypeBytes, endian)
		if err != nil {
			return err
		}
		tag.DataType = dataType

		// set count of data
		dataCountBytes, err := br.ReadBytes(4)
		if err != nil {
			return fmt.Errorf("Failed to read count of data in bytes\n")
		}
		dataCount := endian.Uint32(dataCountBytes)
		tag.DataCount = dataCount

		// check data size
		dataTypeSize, err := dataType.ByteSize()
		if err != nil {
			return err
		}
		totalTagDataSize := uint32(dataTypeSize) * tag.DataCount

		dataOrOffset, err := br.ReadBytes(4)
		if err != nil {
			return fmt.Errorf("failed to read dataOrOffset")
		}

		if totalTagDataSize > 4 {
			fmt.Printf("Data > 4 bytes\n")
			// set absolute data offset
			offset := endian.Uint32(dataOrOffset)
			absDataOffset := tiffHeaderStart + int64(offset)

			// save current pos
			currentPos, err := br.Seek(0, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("Failed to save current position\n")
			}

			// jump to offset
			_, err = br.Seek(absDataOffset, io.SeekStart)
			if err != nil {
				return fmt.Errorf("Failed to seek to app1 data tag offset\n")
			}

			// read data in bytes
			dataInBytes, err := br.ReadBytes(int(totalTagDataSize))
			if err != nil {
				return fmt.Errorf("Failed to read dataInBytes\n")
			}

			// decode data
			dataValue, err := metadata.DecodeTagData(dataInBytes, dataType, dataCount, endian)
			//fmt.Printf("DV Bytes: %x\n", dataValue)
			tag.Data = dataValue

			fmt.Println("Name: ", tag.Name)
			fmt.Println("Value: ", tag.DataString())

			// set reader back to original pos
			_, err = br.Seek(currentPos, io.SeekStart)
			if err != nil {
				return fmt.Errorf("Failed to seek to app1 data tag offset\n")
			}
		} else {
			fmt.Printf("Data < 4 bytes")
			actualInLineData := dataOrOffset[:totalTagDataSize]

			// Decode in line data
			decodedData, err := metadata.DecodeTagData(actualInLineData, tag.DataType, tag.DataCount, endian)
			if err != nil {
				return fmt.Errorf("failed to decode inline data for tag 0x%04X: %w", tag.ID, err)
			}
			tag.Data = decodedData
		}
	}
	print("END OF PARSE ID")
	return nil
}
