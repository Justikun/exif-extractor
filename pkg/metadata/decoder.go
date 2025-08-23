package metadata

import (
	"fmt"
	"encoding/binary"
)

func GetDataTypeFromBytes(b []byte, byteOrder binary.ByteOrder) (DataType, error) {
	dataValue := int(byteOrder.Uint16(b))

	switch DataType(dataValue) { // Cast the int to DataType for the switch
	case TypeByte:
		return TypeByte, nil
	case TypeAscii:
		return TypeAscii, nil
	case TypeShort:
		return TypeShort, nil
	case TypeLong:
		return TypeLong, nil
	case TypeRational:
		return TypeRational, nil
	case TypeSByte:
		return TypeSByte, nil
	case TypeUndefined:
		return TypeUndefined, nil
	case TypeSShort:
		return TypeSShort, nil
	case TypeSLong:
		return TypeSLong, nil
	case TypeSRational:
		return TypeSRational, nil
	case TypeFloat:
		return TypeFloat, nil
	case TypeDouble:
		return TypeDouble, nil
	default:
		return 0, fmt.Errorf("Failed to GetDataTypeBytes. Unknown data type value: %d", dataValue)
	}
}

func DecodeTagByteData(data []byte, endian binary.ByteOrder) {
	dataIntValue := int(binary.ByteOrder.Uin)

	switch DataType()
	// check data type
	fmt.Printf("%s\n\n", data)
}




