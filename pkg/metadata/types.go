package metadata

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type ImageData struct {
	ImagePath string
	MetaData  MetaData
}

type MetaData struct {
	MainTags   []IFDtag
	ExifTags   []IFDtag
	IntropTags []IFDtag
	GPStags    []IFDtag
}

type IFDtag struct {
	ID        uint16
	Name      string
	DataType  DataType
	DataCount uint32
	Data      any
}

type Rational struct {
	Numerator   uint32
	Denominator uint32
}

type Srational struct {
	Numerator   int32
	Denominator int32
}

type IFDtype string

const (
	IFDMAIN   IFDtype = "main"
	IFDEXIF   IFDtype = "exif"
	IFDINTROP IFDtype = "introp"
	IFDGPS    IFDtype = "gps"
)

type DataType uint16

func (t IFDtag) DataString() string {
	switch v := t.Data.(type) {
	case string:
		return v

	// handle single numerical values
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64:
		return fmt.Sprintf("%v", v)
	case float32, float64:
		return fmt.Sprintf("%v", v)

	// Handle slices of numerical values
	case []uint8:
		return formatSlice("%v", v)
	case []uint16:
		return formatSlice("%v", v)
	case []uint32:
		return formatSlice("%v", v)
	case []uint64:
		return formatSlice("%v", v)
	case []int8:
		return formatSlice("%v", v)
	case []int16:
		return formatSlice("%v", v)
	case []int32:
		return formatSlice("%v", v)
	case []int64:
		return formatSlice("%v", v)
	case []float32:
		return formatSlice("%v", v)
	case []float64:
		return formatSlice("%v", v)

	// Handle rationals
	case Rational:
		return fmt.Sprintf("%d/%d", v.Numerator, v.Denominator)
	case Srational:
		return fmt.Sprintf("%d/%d", v.Numerator, v.Denominator)
	case []Rational:
		return formatRationalSlice(v)
	case []Srational:
		return formatSRationalSlice(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// //////////////////////////////////////
// helper functions for reading tag data
// //////////////////////////////////////
func formatSlice[T any](format string, slice []T) string {
	if len(slice) == 0 {
		return "N/A"
	}

	parts := make([]string, len(slice))

	for i, v := range slice {
		parts[i] = fmt.Sprintf(format, v)
	}
	return fmt.Sprintf("%s", strings.Join(parts, " "))
}

func formatRationalSlice(slice []Rational) string {
	if len(slice) == 0 {
		return "N/A"
	}
	parts := make([]string, len(slice))
	for i, v := range slice {
		parts[i] = fmt.Sprintf("%d/%d", v.Numerator, v.Denominator)
	}
	return fmt.Sprintf("%s", strings.Join(parts, " "))
}

func formatSRationalSlice(slice []Srational) string {
	if len(slice) == 0 {
		return "N/A"
	}
	parts := make([]string, len(slice))
	for i, v := range slice {
		parts[i] = fmt.Sprintf("%d/%d", v.Numerator, v.Denominator)
	}
	return fmt.Sprintf("%s", strings.Join(parts, " "))
}

// //////////////////////////////////////
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

func GetDataTypeString(b []byte, byteOrder binary.ByteOrder) (string, error) {
	dataValue := int(byteOrder.Uint16(b))

	switch DataType(dataValue) {
	case TypeByte:
		return "Byte", nil
	case TypeAscii:
		return "Ascii", nil
	case TypeShort:
		return "Short", nil
	case TypeLong:
		return "Long", nil
	case TypeRational:
		return "Rational", nil
	case TypeSByte:
		return "SByte", nil
	case TypeUndefined:
		return "Undefined", nil
	case TypeSShort:
		return "SShort", nil
	case TypeSLong:
		return "SLong", nil
	case TypeSRational:
		return "SRational", nil
	case TypeFloat:
		return "Float", nil
	case TypeDouble:
		return "Double", nil
	default:
		return "", fmt.Errorf("Failed to get data type string. Unknown data type value: %d", dataValue)
	}
}

func GetDataType(b []byte, byteOrder binary.ByteOrder) (DataType, error) {
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

var typeSize = map[DataType]int{
	TypeByte:      1, // 8 bit unsigned int (1 byte)
	TypeAscii:     1, // 1 byte per char. null-terminated
	TypeShort:     2, // Unsigned 16-bit int (2 byte)
	TypeLong:      4, // Unsigned 32-bit int (4 bytes)
	TypeRational:  8, // Two Long types. numerator(4 bytes), denominator(4 bytes)
	TypeSByte:     1, // 8 bit signed in (1 byte)
	TypeUndefined: 1, // Untyped 8-bit data (raw byte, may be binary or ascii)
	TypeSShort:    2, // 16 bit signed int (2 bytes)
	TypeSLong:     4, // 32 bit signed Int (4 bytes)
	TypeSRational: 8, // Two SLong type. numerator(4 bytes), denominator(4 bytes)
	TypeFloat:     4, // 32 bit IEEE floating point (4 bytes)
	TypeDouble:    8, // 64 bit IEEE floating point (8 bytes)
}

func (dt DataType) String() string {
	switch DataType(dt) { // Cast the int to DataType for the switch
	case TypeByte:
		return "BYTE"
	case TypeAscii:
		return "ASCII"
	case TypeShort:
		return "SHORT"
	case TypeLong:
		return "LONG"
	case TypeRational:
		return "RATIONAL"
	case TypeSByte:
		return "SBYTE"
	case TypeUndefined:
		return "UNDEFINED"
	case TypeSShort:
		return "SSHORT"
	case TypeSLong:
		return "SLONG"
	case TypeSRational:
		return "SRATIONAL"
	case TypeFloat:
		return "FLOAT"
	case TypeDouble:
		return "DOUBLE"
	default:
		return ""
	}
}

func (dt DataType) ByteSize() (int, error) {
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
