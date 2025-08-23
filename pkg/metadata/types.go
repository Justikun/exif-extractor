package metadata

import (
	"fmt"
)

type ImageData struct {
	ImagePath	string
	MetaData	MetaData
}

type MetaData struct {
	Tags		[]IFDTag
	ExifTags	[]IFDTag
	IntropTags	[]IFDTag
	GPSTags		[]IFDTag
}

type IFDTag struct {
	Id			uint16
	Name		string
	DataType	DataType
	Count		uint32
	Data		any
}

type Rational struct {
	Numerator   uint32
	Denominator uint32
}

type SRational struct {
	Numerator   int32
	Denominator int32
}

type IFDType string

const (
	IFDMAIN 		IFDType = "main"
	IFDEXIF			IFDType = "exif"
	IFDINTROP		IFDType = "introp"
	IFDGPS			IFDType = "gps"
)

type DataType uint16

const (
	TypeByte      DataType = 1  // Unsigned 8-bit integer
	TypeAscii     DataType = 2  // 8-bit ASCII character
	TypeShort     DataType = 3  // Unsigned 16-bit integer
	TypeLong      DataType = 4  // Unsigned 32-bit integer
	TypeRational  DataType = 5  // Two LONGs (numerator, denominator)
	TypeSByte     DataType = 6  // Signed 8-bit integer
	TypeUndefined DataType = 7  // 8-bit untyped data
	TypeSShort    DataType = 8  // Signed 16-bit integer
	TypeSLong     DataType = 9  // Signed 32-bit integer
	TypeSRational DataType = 10 // Two SLONGs (signed numerator, denominator)
	TypeFloat     DataType = 11 // 32-bit IEEE floating point
	TypeDouble    DataType = 12 // 64-bit IEEE floating point
)

var TypeSize = map[DataType]uint32 {
	TypeByte:      1,
	TypeAscii:     1,
	TypeShort:     2,
	TypeLong:      4,
	TypeRational:  8,
	TypeSByte:     1,
	TypeUndefined: 1,
	TypeSShort:    2,
	TypeSLong:     4,
	TypeSRational: 8,
	TypeFloat:     4,
	TypeDouble:    8,
}

func (dt DataType) GetDataTypeString() (string, error) {
	switch DataType(dt) { // Cast the int to DataType for the switch
	case TypeByte:
		return "BYTE", nil
	case TypeAscii:
		return "ASCII", nil
	case TypeShort:
		return "SHORT", nil
	case TypeLong:
		return "LONG", nil
	case TypeRational:
		return "RATIONAL", nil
	case TypeSByte:
		return "SBYTE", nil
	case TypeUndefined:
		return "UNDEFINED", nil
	case TypeSShort:
		return "SSHORT", nil
	case TypeSLong:
		return "SLONG", nil
	case TypeSRational:
		return "SRATIONAL", nil
	case TypeFloat:
		return "FLOAT", nil
	case TypeDouble:
		return "DOUBLE", nil
	default:
		return "", fmt.Errorf("Failed to GetDataTypeString. Unknown data type value: %d", dt)
	}
}

func (dt DataType) SizeInBytes() (int, error){
	switch dt {
	case TypeByte, TypeAscii, TypeSByte, TypeUndefined:
		return 1, nil
	case TypeShort, TypeSShort:
		return 2, nil
	case TypeLong, TypeSLong, TypeFloat:
		return 4, nil
	case TypeRational, TypeSRational, TypeDouble:
		return 8, nil
	default:
		return 0, fmt.Errorf("Unkown or invalid datatype for size calculation: %d", dt)
	}
}

func GetNameFromID(id uint16) (string, error) {
	if name, exists := ifdMainTagList[id]; exists {
		return name, nil
	}
	return "", fmt.Errorf("No Tag ID found")
}
